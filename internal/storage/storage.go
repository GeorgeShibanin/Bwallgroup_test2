package storage

import (
	"context"
	"errors"
)

var (
	StorageError = errors.New("storage")
)

type Balance int
type Client int

type User struct {
	Id      Client
	Balance Balance
}

type Storage interface {
	GetBalance(ctx context.Context, client Client) (Balance, error)
	PutNewUser(ctx context.Context, client Client, balance Balance) (User, error)
}
