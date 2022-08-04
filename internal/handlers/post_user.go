package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/GeorgeShibanin/Bwallgroup_test2/internal/storage"
	"github.com/pkg/errors"
)

func (h *HTTPHandler) HandlePostUser(rw http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	var data HandlerNameRequest
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	if data.Balance < 0 {
		http.Error(rw, "wrong balance", http.StatusBadRequest)
		return
	}
	putErr := h.storage.PutNewUser(ctx, data.User, data.Balance)
	if putErr != nil && !errors.Is(putErr, storage.ErrAlreadyExist) {
		putErr = errors.Wrap(putErr, "can't put user")
		http.Error(rw, putErr.Error(), http.StatusInternalServerError)
	}

	response := HandlerNameResposne{
		User:    data.User,
		Balance: data.Balance,
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
