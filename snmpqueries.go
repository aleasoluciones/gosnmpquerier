package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/eferro/go-snmpqueries/pkg/snmpquery"
)

func generateRandomQueries() <-chan snmpquery.Query {
	out := make(chan snmpquery.Query, 100)
	go func() {
		queryId := 0
		for {
			query := snmpquery.Query{
				Id:          queryId,
				Query:       "Fake query " + strconv.Itoa(queryId),
				Destination: "Fake destination " + strconv.Itoa(queryId),
			}
			out <- query
			queryId += 1
			time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
		}
	}()
	return out
}

func main() {

	input := generateRandomQueries()

	processed := make(chan snmpquery.QueryResponse, 100)
	go func() {
		for query := range input {
			processed <- snmpquery.HandleQuery(query)
		}
	}()

	for response := range processed {
		fmt.Println(response.Query, response.Response)
	}
}
