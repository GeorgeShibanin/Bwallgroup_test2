package handlers

import (
	"encoding/json"
	"github.com/GeorgeShibanin/Bwallgroup_test2/internal/storage"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
	"strings"
)

func (h *HTTPHandler) HandleGetUrl(rw http.ResponseWriter, r *http.Request) {

	user_id, err := strconv.Atoi(
		strings.Trim(r.URL.Path, "/"),
	)
	if err != nil {

	}
	balance, err := h.storage.GetBalance(r.Context(), storage.Client(user_id))

	if err != nil {
		http.NotFound(rw, r)
		return
	}
	//http.Redirect(rw, r, string(url), http.StatusPermanentRedirect)
	response := PutResponseData{
		Result: balance,
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
