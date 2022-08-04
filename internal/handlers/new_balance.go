package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/GeorgeShibanin/Bwallgroup_test2/internal/broker"
	"github.com/pkg/errors"
)

func (h *HTTPHandler) HandleNewBalance(rw http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	var data HandlerNameResposne
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	trxID, err := h.storage.CreateTransaction(ctx, data.User, data.Balance)
	if err != nil {
		err = errors.Wrap(err, "can't create trx")
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	h.broker.ApplyTransaction(broker.Transaction{
		ID:       trxID,
		ClientID: data.User,
		Amount:   data.Balance,
	})

	response := ResponseTrx{
		User:          data.User,
		TransactionID: trxID,
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
