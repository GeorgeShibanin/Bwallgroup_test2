package handlers

import (
	"encoding/json"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
	"strings"
)

func (h *HTTPHandler) HandleGetBalance(rw http.ResponseWriter, r *http.Request) {
	clientID, err := strconv.Atoi(
		strings.Trim(r.URL.Path, "/"),
	)
	if err != nil {
		http.NotFound(rw, r)
		return
	}

	balance, err := h.storage.GetBalance(r.Context(), int64(clientID))
	if err != nil {
		http.NotFound(rw, r)
		return
	}

	response := PutResponseData{
		Balance: strconv.FormatInt(balance, 10),
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
