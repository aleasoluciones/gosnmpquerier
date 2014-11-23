// Copyright 2014 The GoSNMPQuerier Authors. All rights reserved.  Use of this
// source code is governed by a MIT-style license that can be found in the
// LICENSE file.

package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/aleasoluciones/gosnmpquerier"
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
	close(input)
}

func printResults(processed chan gosnmpquerier.Query) {
	for query := range processed {
		fmt.Println("Result", query)
	}
}

func main() {
	querier := gosnmpquerier.NewAsyncQuerier(CONTENTION, 3, 3*time.Second)

	go readQueriesFromStdin(querier.Input)

	printResults(querier.Output)
}
