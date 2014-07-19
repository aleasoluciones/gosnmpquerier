package snmpquery

import (
	"time"

	"github.com/soniah/gosnmp"
)

type OpSnmp int32

const (
	GET  = 0
	WALK = 1
)

type Query struct {
	Id          int
	Cmd         OpSnmp
	Community   string
	Oid         string
	Timeout     time.Duration
	Retries     int
	Destination string
	Response    []gosnmp.SnmpPDU
	Error       error
}

func Process(input chan Query, processed chan Query, conntention int) {
	m := make(map[string]chan Query)

	for query := range input {
		_, exists := m[query.Destination]
		if exists == false {
			channel_tmp := make(chan Query, 10)
			m[query.Destination] = channel_tmp
			for i := 0; i < conntention; i++ {
				go processQueriesFromChannel(channel_tmp, processed)
			}
		}
		m[query.Destination] <- query
	}
}

func handleQuery(query *Query) {
	switch query.Cmd {
	case WALK:
		query.Response, query.Error = walk(query.Destination, query.Community, query.Oid, query.Timeout, query.Retries)
	case GET:
		query.Response, query.Error = get(query.Destination, query.Community, query.Oid, query.Timeout, query.Retries)
	}
}

func processQueriesFromChannel(input chan Query, processed chan Query) {
	for query := range input {
		handleQuery(&query)
		processed <- query
	}
}

type QueryWithOutputChannel struct {
	query           Query
	responseChannel chan Query
}

type SynchronousQuerier struct {
	Input chan QueryWithOutputChannel
}

func NewSynchronousQuerier(contention int) *SynchronousQuerier {
	querier := SynchronousQuerier{
		Input: make(chan QueryWithOutputChannel),
	}

	go ProcessAndDispatchQueries(querier.Input, contention)
	return &querier
}

func (querier *SynchronousQuerier) ExecuteQuery(query Query) Query {
	output := make(chan Query)
	querier.Input <- QueryWithOutputChannel{query, output}
	processedQuery := <-output
	return processedQuery
}

func ProcessAndDispatchQueries(input chan QueryWithOutputChannel, contention int) {
	inputQueries := make(chan Query, 10)
	processed := make(chan Query, 10)

	go Process(inputQueries, processed, contention)

	m := make(map[int]chan Query)
	i := 0
	for {
		select {
		case queryWithOutputChannel := <-input:
			queryWithOutputChannel.query.Id = i
			i += 1
			m[queryWithOutputChannel.query.Id] = queryWithOutputChannel.responseChannel
			inputQueries <- queryWithOutputChannel.query
		case processedQuery := <-processed:
			m[processedQuery.Id] <- processedQuery
			delete(m, processedQuery.Id)
		}
	}
}
