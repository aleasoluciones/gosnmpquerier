package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/eferro/go-snmpqueries/pkg/snmpquery"
)

const (
	CONTENTION = 4
)

type QueryMessage struct {
	Command        string
	Destination    string
	Community      string
	Oid            string
	Timeout        int
	Retries        int
	AdditionalInfo interface{}
}

func queryFromJson(jsonText string, queryId int) *snmpquery.Query {
	var m QueryMessage
	m.Timeout = 2
	m.Retries = 1

	b := []byte(jsonText)
	err := json.Unmarshal(b, &m)
	if err != nil {
		fmt.Println("Invalid jsonText format", err, jsonText)
		return nil
	}

	cmd, err := convertCommand(m.Command)
	return &snmpquery.Query{
		Id:          queryId,
		Cmd:         cmd,
		Community:   m.Community,
		Oid:         m.Oid,
		Destination: m.Destination,
		Timeout:     time.Duration(m.Timeout) * time.Second,
		Retries:     m.Retries,
	}
}

func convertCommand(command string) (snmpquery.OpSnmp, error) {
	switch command {
	case "walk":
		return snmpquery.WALK, nil
	case "get":
		return snmpquery.GET, nil
	default:
		return 0, fmt.Errorf("Unsupported command %s ", command)
	}
}

type QueryWithOutputChannel struct {
	query           snmpquery.Query
	responseChannel chan snmpquery.Query
}

func ProcessSynchronous(input chan QueryWithOutputChannel) {

	inputQueries := make(chan snmpquery.Query, 10)
	processed := make(chan snmpquery.Query, 10)

	go snmpquery.Process(inputQueries, processed, CONTENTION)

	m := make(map[int]chan snmpquery.Query)
	for {
		select {
		case queryWithOutputChannel := <-input:
			m[queryWithOutputChannel.query.Id] = queryWithOutputChannel.responseChannel
			inputQueries <- queryWithOutputChannel.query
		case processedQuery := <-processed:
			m[processedQuery.Id] <- processedQuery
			// borrar la entrada del mapa....
		}
	}
}

func executeQuery(queryChannel chan QueryWithOutputChannel, query snmpquery.Query) snmpquery.Query {
	output := make(chan snmpquery.Query)
	queryChannel <- QueryWithOutputChannel{query, output}
	processedQuery := <-output
	return processedQuery
}

func main() {
	var input = make(chan QueryWithOutputChannel)
	go ProcessSynchronous(input)

	for id := 0; id < 10; id++ {
		q := queryFromJson(`{"Command":"get", "Destination":"localhost", "Community":"public", "Oid":"1.3.6.1.2.1.31.1.1.1.6.1"}`, id)
		fmt.Println("Result:", executeQuery(input, *q))
	}
}
