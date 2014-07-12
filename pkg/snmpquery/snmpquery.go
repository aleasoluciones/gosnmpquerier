package snmpquery

import (
	"math/rand"
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

func HandleQuery(query Query) QueryResponse {
	time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
	return QueryResponse{
		Id:       query.Id,
		Response: "whatever",
		Query:    query,
	}
}
