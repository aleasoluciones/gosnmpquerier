package snmpquery

import (
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

func Process(input chan Query, processed chan Query) {
	for query := range input {

		HandleQuery(&query)
		processed <- query
	}
}
