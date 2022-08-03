package handlers

import (
	"github.com/GeorgeShibanin/Bwallgroup_test2/internal/storage"
)

type HTTPHandler struct {
	storage storage.Storage
}

func NewHTTPHandler(storage storage.Storage) *HTTPHandler {
	return &HTTPHandler{
		storage: storage,
	}
}

type PutResponseData struct {
	Result string `json:"result"`
}

type PutResponseNewUser struct {
	User    int `json:"user_id"`
	Balance int `json:"balance"`
}
