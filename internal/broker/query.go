package broker

type Query struct {
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
	q.active = true
}

func (q *Query) Deactivate() {
	q.active = false
}

func (q *Query) GetChannel() chan Transaction {
	return q.query
}

func (q *Query) Add(tr Transaction) {
	q.query <- tr
}
