package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

type Query struct {
	Query       string
	Destination string
}

type QueryResponse struct {
	Response string
	Query    Query
}

func generateRandomQueries() <-chan Query {
	out := make(chan Query)
	go func() {
		queryId := 0
		for {
			query := Query{
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
	for query := range input {
		fmt.Println(query.Query)
	}
}
