package snmpquery

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

type Query struct {
	Id          int
	Query       string
	Destination string
	Response    string
	Error       int
}

func HandleQuery(query *Query) {
	time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
	query.Response = "whatever " + strconv.Itoa(rand.Intn(1e3))
}

func ProcessQueriesFromChannel(input chan Query, processed chan Query) {
	for query := range input {
		fmt.Println("Begin processing", query)
		HandleQuery(&query)
		fmt.Println("Processed", query)
		processed <- query
	}
}

func Process(input chan Query, processed chan Query) {
	m := make(map[string]chan Query)

	for query := range input {
		fmt.Println("EFA processing query ", query.Destination, query)

		channel, exists := m[query.Destination]
		if exists == false {
			channel := make(chan Query, 10)
			m[query.Destination] = channel
			fmt.Println("EFA creating a go routine for ", query.Destination)
			go ProcessQueriesFromChannel(channel, processed)
		}
		channel <- query

	}
}
