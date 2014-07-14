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

func Process(input chan Query, processed chan Query, conntention int) {
	m := make(map[string]chan Query)

	for query := range input {
		_, exists := m[query.Destination]
		if exists == false {
			channel_tmp := make(chan Query, 10)
			m[query.Destination] = channel_tmp
			for i := 0; i < conntention; i++ {
				go processQueriesFromChannel(channel_tmp, processed)
			}
		}
		m[query.Destination] <- query
	}
}

func handleQuery(query *Query) {
	time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
	query.Response = "whatever " + strconv.Itoa(rand.Intn(1e3))
}

func processQueriesFromChannel(input chan Query, processed chan Query) {
	for query := range input {
		handleQuery(&query)
		processed <- query
	}
}
