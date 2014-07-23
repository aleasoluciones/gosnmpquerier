package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/eferro/gosnmpquerier"
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

func readQueriesFromStdin(input chan gosnmpquerier.Query) {

	inputLines := make(chan string, 10)
	go readLinesFromStdin(inputLines)

	queryId := 0
	for line := range inputLines {
		query, err := gosnmpquerier.FromJson(line)
		if err != nil {
			fmt.Println("Invalid line:", line, err)
		} else {
			query.Id = queryId
			input <- *query
			queryId += 1
		}
	}
}

func printResults(processed chan gosnmpquerier.Query) {
	for query := range processed {
		fmt.Println("Result", query)
	}
}

func main() {

	querier := gosnmpquerier.NewAsyncQuerier(CONTENTION)

	go readQueriesFromStdin(querier.Input)

	printResults(querier.Output)
}
