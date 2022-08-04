package handlers

import (
	"github.com/GeorgeShibanin/Bwallgroup_test2/internal/broker"
	"github.com/GeorgeShibanin/Bwallgroup_test2/internal/storage"
)

type HTTPHandler struct {
	storage storage.Storage
	broker  *broker.Broker
}

func NewHTTPHandler(st storage.Storage, br *broker.Broker) *HTTPHandler {
	return &HTTPHandler{
		storage: st,
		broker:  br,
	}
}

type PutResponseData struct {
	Balance string `json:"balance"`
}

type HandlerNameRequest struct {
	User    int64 `json:"user_id"`
	Balance int64 `json:"balance"`
}

type HandlerNameResposne struct {
	User    int64 `json:"user_id"`
	Balance int64 `json:"balance"`
}

type ResponseTrx struct {
	User          int64 `json:"user_id"`
	TransactionID int64 `json:"transactionID"`
}
