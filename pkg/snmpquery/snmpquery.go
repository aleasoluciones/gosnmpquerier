package snmpquery

import (
	"time"

	"github.com/eferro/gosnmp"
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

type QueryWithOutputChannel struct {
	query           Query
	responseChannel chan Query
}

type SyncQuerier struct {
	Input        chan QueryWithOutputChannel
	asyncQuerier *AsyncQuerier
}

func NewSyncQuerier(contention int) *SyncQuerier {
	querier := SyncQuerier{
		Input:        make(chan QueryWithOutputChannel),
		asyncQuerier: NewAsyncQuerier(contention),
	}
	return &querier
}

func (querier *SyncQuerier) ExecuteQuery(query Query) Query {
	output := make(chan Query)
	querier.Input <- QueryWithOutputChannel{query, output}
	processedQuery := <-output
	return processedQuery
}

func processAndDispatchQueries(input chan QueryWithOutputChannel, contention int) {
	// inputQueries := make(chan Query, 10)
	// processed := make(chan Query, 10)

	// asyncQuerier := NewAsyncQuerier(contention)
	// asyncQuerier.process()
	// go process(inputQueries, processed, contention)

	// m := make(map[int]chan Query)
	// i := 0
	// for {
	// 	select {
	// 	case queryWithOutputChannel := <-input:
	// 		queryWithOutputChannel.query.Id = i
	// 		i += 1
	// 		m[queryWithOutputChannel.query.Id] = queryWithOutputChannel.responseChannel
	// 		inputQueries <- queryWithOutputChannel.query
	// 	case processedQuery := <-processed:
	// 		m[processedQuery.Id] <- processedQuery
	// 		delete(m, processedQuery.Id)
	// 	}
	// }
}
