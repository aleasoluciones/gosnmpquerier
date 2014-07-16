package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/eferro/go-snmpqueries/pkg/snmpquery"
)

const (
	CONTENTION = 4
)

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
		query, err := snmpquery.FromJson(line)
		if err != nil {
			fmt.Println("Invalid line:", line, err)
		} else {
			query.Id = queryId
			input <- *query
			queryId += 1
		}
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
