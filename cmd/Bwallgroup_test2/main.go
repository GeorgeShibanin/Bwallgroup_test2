package main

import (
	"context"
	"fmt"
	"github.com/GeorgeShibanin/Bwallgroup_test2/internal/broker"
	"github.com/GeorgeShibanin/Bwallgroup_test2/internal/config"
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
	trxBroker := broker.InitBroker(store)
	err := initBroker(context.Background(), store, trxBroker)
	if err != nil {
		log.Println(err)
		return nil
	}

	handler := handlers.NewHTTPHandler(store, trxBroker)
	r.HandleFunc("/{user_id:\\w{1}}", handler.HandleGetBalance).Methods(http.MethodGet)
	r.HandleFunc("/", handler.HandlePostUser)
	r.HandleFunc("/trx", handler.HandleNewBalance)

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

func initBroker(ctx context.Context, st storage.Storage, br *broker.Broker) error {
	trxs, err := st.GetActiveTransactions(ctx)
	if err != nil {
		return fmt.Errorf("cant get Active transactions - %w", err)
	}
	for _, trx := range trxs {
		br.ApplyTransaction(broker.Transaction{
			ID:       trx.Id,
			ClientID: trx.IdClient,
			Amount:   trx.OperationSum,
		})
	}
	return nil
}
