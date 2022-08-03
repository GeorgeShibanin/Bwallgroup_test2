package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"

	"github.com/GeorgeShibanin/Bwallgroup_test2/internal/storage"
	"github.com/pkg/errors"
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

func (s *StoragePostgres) PutNewUser(ctx context.Context, key storage.ShortedURL, url storage.URL) (storage.ShortedURL, error) {

}

func (s *StoragePostgres) GetBalance(ctx context.Context, user_id storage.Client) (storage.Balance, error) {
}

func (s *StoragePostgres) PatchUserBalance(ctx context.Context, client Client, balance Balance) (Balance, error) {

}
