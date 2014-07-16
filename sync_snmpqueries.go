package main

import (
	"fmt"

	"github.com/eferro/go-snmpqueries/pkg/snmpquery"
)

const (
	CONTENTION = 4
)

func main() {
	var input = make(chan snmpquery.QueryWithOutputChannel)
	go snmpquery.ProcessAndDispatchQueries(input, CONTENTION)

	for id := 0; id < 10; id++ {
		q, _ := snmpquery.FromJson(`{"command":"walk", "destination":"localhost", "community":"public", "oid":"1.3.6.1"}`)
		q.Id = id

		fmt.Println("Result:", snmpquery.ExecuteQuery(input, *q))
	}
}
