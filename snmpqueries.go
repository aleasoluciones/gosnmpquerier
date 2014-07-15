package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
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

func readLinesFromStdin(inputLines chan string) {
	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			close(inputLines)
			return
		}
		inputLines <- line
	}
}

func readQueriesFromStdin(input chan snmpquery.Query) {

	inputLines := make(chan string, 10)
	go readLinesFromStdin(inputLines)

	queryId := 0
	for line := range inputLines {

		var m QueryMessage
		m.Timeout = 2
		m.Retries = 1

		b := []byte(line)
		err := json.Unmarshal(b, &m)
		if err != nil {
			fmt.Println("Invalid line format", err, line)
		}

		cmd, err := convertCommand(m.Command)
		query := snmpquery.Query{
			Id:          queryId,
			Cmd:         cmd,
			Community:   m.Community,
			Oid:         m.Oid,
			Destination: m.Destination,
			Timeout:     time.Duration(m.Timeout) * time.Second,
			Retries:     m.Retries,
		}
		fmt.Println("AÑADIMOS QUERY", query.Id)
		input <- query
		queryId += 1
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

func printResults(processed chan snmpquery.Query) {
	for query := range processed {
		fmt.Println("Result", query)
	}
}

func main() {
	input := make(chan snmpquery.Query, 10)
	processed := make(chan snmpquery.Query, 10)

	go readQueriesFromStdin(input)
	go snmpquery.Process(input, processed, CONTENTION)

	printResults(processed)
}
