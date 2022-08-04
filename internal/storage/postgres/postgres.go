package postgres

import (
	"context"
	"fmt"
	"github.com/GeorgeShibanin/Bwallgroup_test2/internal/storage"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"log"
	"strconv"
	"time"
)

const (
	GetUserByIDQuery          = `SELECT id, balance FROM client WHERE id = $1`
	GetActiveTransactionQuery = `SELECT * FROM query WHERE operation_accepted = false ORDER BY created_at`
	GetTxByIDQuery            = `SELECT id, client_id, operation_sum, operation_accepted, created_at FROM query WHERE id = $1 `

	InsertUserQuery        = `INSERT INTO client (id, balance) values ($1, $2)`
	InsertTransactionQuery = `INSERT INTO query (client_id, operation_sum, operation_accepted, created_at) 
								values($1, $2, $3, $4) RETURNING id`

	UpdateStatusQuery = `UPDATE query SET operation_accepted = $1 WHERE id = $2
								RETURNING id`
	UpdateBalanceQuery = `UPDATE client SET balance = $1 WHERE id = $2
								RETURNING id`

	dsnTemplate = "postgres://%s:%s@%s:%v/%s"
)

// запрос примерно такой
// SELECT * FROM query WHERE operation_accepted = false ORDER BY created_at

type StoragePostgres struct {
	conn *pgx.Conn
}

func initConnection(conn *pgx.Conn) *StoragePostgres {
	return &StoragePostgres{conn: conn}
}

func Init(ctx context.Context, host, user, db, password string, port uint16) (*StoragePostgres, error) {
	//подключение к базе через переменные окружения
	conn, err := pgx.Connect(ctx, fmt.Sprintf(dsnTemplate, user, password, host, port, db))
	if err != nil {
		return nil, errors.Wrap(err, "can't connect to postgres")
	}
	return initConnection(conn), nil
}

func (s *StoragePostgres) PutNewUser(ctx context.Context, clientID, amount int64) error {
	tx, err := s.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return errors.Wrap(err, "can't create tx")
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()
	user := &Client{}
	err = tx.QueryRow(ctx, GetUserByIDQuery, clientID).Scan(&user.Id, &user.Balace)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return errors.Wrap(err, "can't get by id")
	}
	if user.Id != 0 {
		return storage.ErrAlreadyExist
	}

	tag, err := tx.Exec(ctx, InsertUserQuery, clientID, amount)
	if err != nil {
		return errors.Wrap(err, "can't insert link")
	}

	if tag.RowsAffected() != 1 {
		return errors.Wrap(err, fmt.Sprintf("unexpected rows affected value: %v", tag.RowsAffected()))
	}

	return nil
}

func (s *StoragePostgres) GetBalance(ctx context.Context, clientID int64) (int64, error) {
	client := &Client{}
	//получаем из базы значение по ключу
	err := s.conn.QueryRow(ctx, GetUserByIDQuery, clientID).
		Scan(&client.Id, &client.Balace)
	if err != nil {
		return 0, fmt.Errorf("something went wrong - %w", storage.StorageError)
	}
	return client.Balace, err
}

func (s *StoragePostgres) PatchUserBalance(ctx context.Context, clientID, transactionID int64) (int64, error) {
	tx, err := s.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return 0, errors.Wrap(err, "can't create tx")
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	user := &Client{}
	//Получаем информацию о клиенте из первой таблицы
	err = s.conn.QueryRow(ctx, GetUserByIDQuery, clientID).
		Scan(&user.Id, &user.Balace)
	if err != nil {
		return 0, fmt.Errorf("something went wrong - %w", storage.StorageError)
	}
	//баланс клиента
	currentBalance := user.Balace

	query := &storage.Transaction{}
	//Получаем транзакцию по Id
	err = tx.QueryRow(ctx, GetTxByIDQuery, strconv.FormatInt(transactionID, 10)).Scan(&query.Id, &query.IdClient,
		&query.OperationSum, &query.OperationAccepted, &query.CreatedAt)
	if err != nil {
		return 0, errors.Wrap(err, "something went wrong with getting transaction")
	}
	//можем ли изменить баланс клиента?
	if currentBalance+query.OperationSum < 0 {
		return 0, fmt.Errorf("not enough balance - %w", storage.StorageError)
	}

	//обновляем статуст транзакции
	tag, err := s.conn.Exec(ctx, UpdateStatusQuery, true, transactionID)
	if err != nil {
		return 0, fmt.Errorf("cant update balance - %w", storage.StorageError)
	}
	if tag.RowsAffected() != 1 {
		return 0, errors.Wrap(err, fmt.Sprintf("unexpected rows affected value: %v", tag.RowsAffected()))
	}

	//обновляем значения для клиента
	tag, err = s.conn.Exec(ctx, UpdateBalanceQuery, currentBalance+query.OperationSum, clientID)
	if err != nil {
		return 0, fmt.Errorf("cant update balance - %w", storage.StorageError)
	}
	if tag.RowsAffected() != 1 {
		return 0, errors.Wrap(err, fmt.Sprintf("unexpected rows affected value: %v", tag.RowsAffected()))
	}

	return query.Id, nil
}

// CreateTransaction - создает транзакцию и возвращает id
func (s *StoragePostgres) CreateTransaction(ctx context.Context, clientID, amount int64) (int64, error) {
	query := &storage.Transaction{}
	err := s.conn.QueryRow(ctx, InsertTransactionQuery, clientID, amount, false, time.Now().UTC().Format(time.RFC3339)).
		Scan(&query.Id)
	if err != nil {
		return 0, fmt.Errorf("cant add new transaction - %w", storage.StorageError)
	}
	return query.Id, nil
}

func (s *StoragePostgres) GetActiveTransactions(ctx context.Context) ([]*storage.Transaction, error) {
	trxs := make([]*storage.Transaction, 0)
	rows, err := s.conn.Query(ctx, GetActiveTransactionQuery)
	if err != nil {
		log.Println("cant get active trx")
		return nil, fmt.Errorf("cant get active transaction %w", storage.StorageError)
	}

	for rows.Next() {
		trx := &storage.Transaction{}
		err = rows.Scan(&trx.Id, &trx.IdClient, &trx.OperationSum, &trx.OperationAccepted, &trx.CreatedAt)
		if err != nil {
			log.Println("wrong scan")
			return nil, fmt.Errorf("cant scan rows in every active transaction - %w", err)
		}
		trxs = append(trxs, trx)
	}

	return trxs, nil
}
