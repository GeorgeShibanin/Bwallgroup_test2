package postgres

import (
	"context"
	"fmt"
	"github.com/GeorgeShibanin/Bwallgroup_test2/internal/storage"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"time"
)

const (
	GetUserByIDQuery = `SELECT id, balance FROM client WHERE id = $1`
	InsertUserQuery  = `INSERT INTO client (id, balance) values ($1, $2)`
	PutSummQuery     = `INSERT INTO query (operation_sum, operation_accepted, created_at) values($1, $2, $3)`
	dsnTemplate      = "postgres://%s:%s@%s:%v/%s"
)

type StoragePostgres struct {
	conn postgresInterface
}

type postgresInterface interface {
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

func initConnection(conn postgresInterface) *StoragePostgres {
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

func (s *StoragePostgres) PutNewUser(ctx context.Context, client storage.Client, balance storage.Balance) error {
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
	err = tx.QueryRow(ctx, GetUserByIDQuery, client).Scan(&user.Id, &user.Balace)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return errors.Wrap(err, "can't get by id")
	}
	if user.Id != 0 {
		return storage.ErrAlreadyExist
	}

	tag, err := tx.Exec(ctx, InsertUserQuery, client, balance)
	if err != nil {
		return errors.Wrap(err, "can't insert link")
	}

	if tag.RowsAffected() != 1 {
		return errors.Wrap(err, fmt.Sprintf("unexpected rows affected value: %v", tag.RowsAffected()))
	}

	return nil
}

func (s *StoragePostgres) GetBalance(ctx context.Context, user_id storage.Client) (storage.Balance, error) {
	client := &Client{}
	//получаем из базы значение по ключу
	err := s.conn.QueryRow(ctx, GetUserByIDQuery, user_id).
		Scan(&client.Id, &client.Balace)
	if err != nil {
		return 0, fmt.Errorf("something went wrong - %w", storage.StorageError)
	}
	return storage.Balance(client.Balace), err
}

func (s *StoragePostgres) PatchUserBalance(ctx context.Context, client storage.Client, balance storage.Balance) (storage.Balance, error) {
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
	//получаем из базы значение по ключу
	err = tx.QueryRow(ctx, GetUserByIDQuery, client).
		Scan(&user.Id, &user.Balace)
	if err != nil {
		return 0, fmt.Errorf("something went wrong - %w", storage.StorageError)
	}

	current_balance := user.Balace
	if current_balance+int(balance) < 0 {
		return 0, fmt.Errorf("balance not enough - %w", storage.StorageError)
	}
	tag, err := tx.Exec(ctx, PutSummQuery, balance, false, time.Now())
	if err != nil {
		return 0, errors.Wrap(err, "can't update balance")
	}

	if tag.RowsAffected() != 1 {
		return 0, errors.Wrap(err, fmt.Sprintf("unexpected rows affected value: %v", tag.RowsAffected()))
	}
	return storage.Balance(current_balance + int(balance)), nil
}
