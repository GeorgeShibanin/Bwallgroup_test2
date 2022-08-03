package handlers

import (
	"context"
	"encoding/json"
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

}
