package storage

import (
	"context"
	"errors"
)

var (
	StorageError    = errors.New("storage")
	ErrAlreadyExist = errors.New("client")
)

type Balance int
type Client int

type Storage interface {
	//Получение баланса пользователя
	GetBalance(ctx context.Context, client Client) (Balance, error)
	//Добавление нового пользователя
	PutNewUser(ctx context.Context, client Client, balance Balance) error
	//Изменение баланса пользователя
	PatchUserBalance(ctx context.Context, client Client, balance Balance) (Balance, error)
}
