**Запуск**
```
docker-compose up --build
```

**Добавить клиента**

```
curl -d '{"user_id": 1, "balance": 126}' -X POST http://localhost:8080/new
```

**Вывести баланс**


```
curl -X GET  http://localhost:8080/1
```

где 1 - user_id

**Транзакция**
user_id и на сколько изменить баланс

```
curl -d '{"user_id": 1, "balance": 126}' -X POST http://localhost:8080/trx
```

**Особенности реализаци**

Существует таблица client(id, balance). Все транзакции сохраняются во второй таблице query.
У каждой транзакции есть статус operation_accepted bool

Для каждого клиента есть очередь в виде канала. Для этого канала мы создаём горутину, которая его читает
```
if !ok || !query.IsActive() {
		go b.initReader(query)
	} 
```
Она либо читает транзакции из канала, либо деактивируется через какое то время
```
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
```
Для работы со статусом транзакции используется mutex
```
func (b *Broker) ApplyTransaction(trx Transaction) {
	fmt.Println("ApplyTransaction: ", trx)
	b.mutex.Lock()
	defer b.mutex.Unlock()
```
