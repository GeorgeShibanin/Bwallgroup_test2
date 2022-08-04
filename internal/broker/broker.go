package broker

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/GeorgeShibanin/Bwallgroup_test2/internal/storage"
)

type Transaction struct {
	ID       int64
	ClientID int64
	Amount   int64
}

type Broker struct {
	mutex   sync.Mutex
	queries map[int64]*Query
	storage storage.Storage
}

func InitBroker(st storage.Storage) *Broker {
	return &Broker{
		mutex:   sync.Mutex{},
		queries: make(map[int64]*Query),
		storage: st,
	}
}

func (b *Broker) ApplyTransaction(trx Transaction) {
	fmt.Println("ApplyTransaction: ", trx)
	b.mutex.Lock()
	defer b.mutex.Unlock()

	query, ok := b.queries[trx.ClientID]
	if !ok || query == nil {
		fmt.Println("query initialization: ", trx)
		query = InitQuery()
		b.queries[trx.ClientID] = query
	}

	if !ok || !query.IsActive() {
		go b.initReader(query)
	}

	query.Add(trx)
}

func (b *Broker) initReader(q *Query) {
	fmt.Println("initReader: ", q)
	q.Activate()
	defer q.Deactivate()
	for {
		select {
		case trx := <-q.GetChannel():
			fmt.Println("update balance: ", trx.ClientID)

			_, err := b.storage.PatchUserBalance(context.Background(), trx.ClientID, trx.ID)
			if err != nil {
				log.Printf("can't process transaction: %s", err.Error())
			}

		case <-time.After(3 * time.Second):
			fmt.Println("deactivate query")
			return
		}
	}
}
