package main

import (
	"context"
	"github.com/GeorgeShibanin/Bwallgroup_test2/internal/handlers"
	"github.com/GeorgeShibanin/Bwallgroup_test2/internal/storage"
	"github.com/GeorgeShibanin/Bwallgroup_test2/internal/storage/postgres"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

func main() {
	srv := NewServer()
	log.Printf("Start serving on %s", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}

func NewServer() *http.Server {
	r := mux.NewRouter()

	var store storage.Storage

	store = initPostgres()

	handler := handlers.NewHTTPHandler(store)
	r.HandleFunc("/{shortUrl:\\w{10}}", handler.HandleGetBalance).Methods(http.MethodGet)
	r.HandleFunc("/", handler.HandleP).Methods(http.MethodGet)

	return &http.Server{
		Handler:      r,
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
}

func initPostgres() *postgres.StoragePostgres {
	store, err := postgres.Init(
		context.Background(),
		config.PostgresHost,
		config.PostgresUser,
		config.PostgresDB,
		config.PostgresPassword,
		config.PostgresPort,
	)
	if err != nil {
		log.Fatalf("can't init postgres connection: %s", err.Error())
	}
	return store
}
