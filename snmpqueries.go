package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/eferro/go-snmpqueries/pkg/snmpquery"
)

func generateRandomQueries(input chan snmpquery.Query) {
	queryId := 0
	for {
		query := snmpquery.Query{
			Id:          queryId,
			Query:       "Fake query " + strconv.Itoa(queryId),
			Destination: "Fake destination " + strconv.Itoa(rand.Intn(10)),
		}
		input <- query
		queryId += 1
		time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
	}
}

func printResults(processed chan snmpquery.Query) {
	for query := range processed {
		fmt.Println(query.Destination, query.Query, query.Response)
	}
}

func main() {
	input := make(chan snmpquery.Query, 10)
	processed := make(chan snmpquery.Query, 10)

	go generateRandomQueries(input)
	go snmpquery.Process(input, processed)

	printResults(processed)
}