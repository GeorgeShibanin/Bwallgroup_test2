package handlers

import (
	"context"
	"encoding/json"
	"github.com/GeorgeShibanin/Bwallgroup_test2/internal/storage"
	"github.com/pkg/errors"
	"net/http"
)

const RetriesCount = 5

func (h *HTTPHandler) HandlePostUrl(rw http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	var data ResponseUser
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	putErr := h.storage.PutNewUser(ctx, storage.Client(data.User), storage.Balance(data.Balance))
	if putErr != nil && !errors.Is(putErr, storage.ErrAlreadyExist) {
		putErr = errors.Wrap(putErr, "can't put url")
		http.Error(rw, putErr.Error(), http.StatusInternalServerError)
	}

	response := ResponseUser{
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
