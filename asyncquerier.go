// Copyright 2014 The GoSNMPQuerier Authors. All rights reserved.  Use of this
// source code is governed by a MIT-style license that can be found in the
// LICENSE file.

package gosnmpquerier

import (
	"log"
	"strconv"
)

type AsyncQuerier struct {
	Input      chan Query
	Output     chan Query
	Contention int
}

func NewAsyncQuerier(contention int) *AsyncQuerier {
	querier := AsyncQuerier{
		Input:      make(chan Query, 10),
		Output:     make(chan Query, 10),
		Contention: contention,
	}
	go querier.process()
	return &querier
}

type destinationProcessorInfo struct {
	input  chan Query
	output chan Query
	done   chan bool
}

func (querier *AsyncQuerier) process() {
	log.Println("AsyncQuerier process begin")
	m := make(map[string]destinationProcessorInfo)

	for query := range querier.Input {
		_, exists := m[query.Destination]
		if exists == false {

			processorInfo := destinationProcessorInfo{
				input:  make(chan Query, 10),
				output: querier.Output,
				done:   make(chan bool, 1),
			}

			m[query.Destination] = processorInfo
			for i := 0; i < querier.Contention; i++ {
				go processQueriesFromChannel(
					processorInfo.input,
					processorInfo.output,
					processorInfo.done,
					string(query.Destination)+strconv.Itoa(i))
			}
		}
		m[query.Destination].input <- query
	}
	log.Println("AsyncQuerier process terminating")

	for destination, processorInfo := range m {
		log.Println("closing:", processorInfo.input)
		close(processorInfo.input)
		for i := 0; i < querier.Contention; i++ {
			<-processorInfo.done
		}
		delete(m, destination)
	}
	log.Println("closing output")
	close(querier.Output)
}

func handleQuery(query *Query) {
	switch query.Cmd {
	case WALK:
		query.Response, query.Error = walk(query.Destination, query.Community, query.Oid, query.Timeout, query.Retries)
	case GET:
		query.Response, query.Error = get(query.Destination, query.Community, query.Oid, query.Timeout, query.Retries)
	}
}

func processQueriesFromChannel(input chan Query, processed chan Query, done chan bool, processorId string) {
	for query := range input {
		handleQuery(&query)
		processed <- query
	}
	done <- true
	log.Println(processorId, "terminated")
}
