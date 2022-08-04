package broker

import "sync"

type Query struct {
	mutex  sync.Mutex
	query  chan Transaction
	active bool
}

func InitQuery() *Query {
	return &Query{query: make(chan Transaction)}
}

func (q *Query) IsActive() bool {

	return q.active
}

func (q *Query) Activate() {
	q.mutex.Lock()
	q.active = true
}

func (q *Query) Deactivate() {
	q.active = false
	q.mutex.Unlock()
}

func (q *Query) GetChannel() chan Transaction {
	return q.query
}

func (q *Query) Add(tr Transaction) {
	q.query <- tr
}
