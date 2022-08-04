package storage

import (
	"context"
	"errors"
	"time"
)

var (
	StorageError    = errors.New("storage")
	ErrAlreadyExist = errors.New("client")
)

type Balance int64
type Client int64

type Transaction struct {
	Id                int64
	IdClient          int64
	OperationSum      int64
	OperationAccepted bool
	CreatedAt         time.Time
}

type Storage interface {
	//Получение баланса пользователя
	GetBalance(ctx context.Context, clientID int64) (int64, error)
	//Добавление нового пользователя
	PutNewUser(ctx context.Context, clientID, amount int64) error
	//Изменение баланса пользователя
	PatchUserBalance(ctx context.Context, clientID, transactionID int64) (int64, error)

	CreateTransaction(ctx context.Context, clientID, amount int64) (int64, error)

	GetActiveTransactions(ctx context.Context) ([]*Transaction, error)
}
