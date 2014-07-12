package main

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
}

type QueryResponse struct {
	Id       int
	Response string
	Query    Query
}

func handleQuery(query Query) QueryResponse {
	time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
	return QueryResponse{
		Id:       query.Id,
		Response: "whatever",
		Query:    query,
	}
}

func generateRandomQueries() <-chan Query {
	out := make(chan Query, 100)
	go func() {
		queryId := 0
		for {
			query := Query{
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

	processed := make(chan QueryResponse, 100)
	go func() {
		for query := range input {
			processed <- handleQuery(query)
		}
	}()

	for response := range processed {
		fmt.Println(response.Query, response.Response)
	}
}
