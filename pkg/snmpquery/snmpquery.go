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
	fmt.Println("Process queries from channel", input)
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

		_, exists := m[query.Destination]
		if exists == false {
			channel_tmp := make(chan Query, 10)
			m[query.Destination] = channel_tmp
			fmt.Println("EFA creating a go routine for ", query.Destination, channel_tmp)
			go ProcessQueriesFromChannel(channel_tmp, processed)
		}

		if len(m) > 10 {
			panic(fmt.Sprintf("Hay más de 10 hilos, ¿qué pasa?"))
		}

		channel, _ := m[query.Destination]
		channel <- query

	}
}
