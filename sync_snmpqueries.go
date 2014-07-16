package main

import (
	"fmt"

	"github.com/eferro/go-snmpqueries/pkg/snmpquery"
)

const (
	CONTENTION = 4
)

type QueryWithOutputChannel struct {
	query           snmpquery.Query
	responseChannel chan snmpquery.Query
}

func ProcessAndDispatchQueries(input chan QueryWithOutputChannel) {

	inputQueries := make(chan snmpquery.Query, 10)
	processed := make(chan snmpquery.Query, 10)

	go snmpquery.Process(inputQueries, processed, CONTENTION)

	m := make(map[int]chan snmpquery.Query)
	for {
		select {
		case queryWithOutputChannel := <-input:
			m[queryWithOutputChannel.query.Id] = queryWithOutputChannel.responseChannel
			inputQueries <- queryWithOutputChannel.query
		case processedQuery := <-processed:
			m[processedQuery.Id] <- processedQuery
			delete(m, processedQuery.Id)
		}
	}
}

func executeQuery(queryChannel chan QueryWithOutputChannel, query snmpquery.Query) snmpquery.Query {
	output := make(chan snmpquery.Query)
	queryChannel <- QueryWithOutputChannel{query, output}
	processedQuery := <-output
	return processedQuery
}

func main() {
	var input = make(chan QueryWithOutputChannel)
	go ProcessAndDispatchQueries(input)

	for id := 0; id < 10; id++ {
		q, _ := snmpquery.FromJson(`{"command":"walk", "destination":"localhost", "community":"public", "oid":"1.3.6.1.2.1.25.1"}`)
		q.Id = id
		fmt.Println("Result:", executeQuery(input, *q))
	}
}
