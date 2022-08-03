package handlers

import (
	"context"
	"encoding/json"
	"github.com/GeorgeShibanin/Bwallgroup_test2/internal/storage"
	"github.com/pkg/errors"
	"net/http"
)

func (h *HTTPHandler) HandleNewBalance(rw http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	var data ResponseUser
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	newBalance, err := h.storage.PatchUserBalance(ctx, storage.Client(data.User),
		storage.Balance(data.Balance))

	response := ResponseUser{
		User:    data.User,
		Balance: int(newBalance),
	}
	rawResponse, err := json.Marshal(response)
	if err != nil {
		err = errors.Wrap(err, "can't marshal response")
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	_, err = rw.Write(rawResponse)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
}
