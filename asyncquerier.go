// Copyright 2014 The GoSNMPQuerier Authors. All rights reserved.  Use of this
// source code is governed by a MIT-style license that can be found in the
// LICENSE file.

package gosnmpquerier

import (
	"fmt"
	"log"
	"strconv"
)

type AsyncQuerier struct {
	Input      chan Query
	Output     chan Query
	Contention int
	SnmpClient SnmpClient
}

func NewAsyncQuerier(contention int) *AsyncQuerier {
	querier := AsyncQuerier{
		Input:      make(chan Query, 10),
		Output:     make(chan Query, 10),
		Contention: contention,
		SnmpClient: newSnmpClient(),
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
			processorInfo := createProcessorInfo(querier.Output)
			m[query.Destination] = processorInfo
			createProcessors(processorInfo, querier, query.Destination)
		}
		m[query.Destination].input <- query
	}
	log.Println("AsyncQuerier process terminating")

	waitUntilProcessorEnd(m, querier.Contention)
	log.Println("closing output")
	close(querier.Output)
}
func createProcessorInfo(output chan Query) destinationProcessorInfo {
	return destinationProcessorInfo{
		input:  make(chan Query, 10),
		output: output,
		done:   make(chan bool, 1),
	}
}

func createProcessors(processorInfo destinationProcessorInfo, AsyncQuerier querier, destination string) {
	for i := 0; i < querier.contention; i++ {
		go processQueriesFromChannel(
			querier,
			processorInfo.input,
			processorInfo.output,
			processorInfo.done,
			string(destination)+string("_")+strconv.Itoa(i))
	}
}

func waitUntilProcessorEnd(m map[string]destinationProcessorInfo, contention int) {
	for destination, processorInfo := range m {
		log.Println("closing:", processorInfo.input)
		close(processorInfo.input)
		for i := 0; i < contention; i++ {
			<-processorInfo.done
		}
		delete(m, destination)
	}
}

func (querier *AsyncQuerier) handleQuery(query *Query) {
	switch query.Cmd {
	case WALK:
		if len(query.Oids) == 1 {
			query.Response, query.Error = querier.SnmpClient.walk(query.Destination, query.Community, query.Oids[0], query.Timeout, query.Retries)
		} else {
			query.Error = fmt.Errorf("Multi Oid Walk not supported")
		}
	case GET:
		query.Response, query.Error = querier.SnmpClient.get(query.Destination, query.Community, query.Oids, query.Timeout, query.Retries)
	case GETNEXT:
		query.Response, query.Error = querier.SnmpClient.getnext(query.Destination, query.Community, query.Oids, query.Timeout, query.Retries)
	}
}

func processQueriesFromChannel(AsyncQuerier querier, input chan Query, processed chan Query, done chan bool, processorId string) {
	for query := range input {
		querier.handleQuery(&query)
		processed <- query
	}
	done <- true
	log.Println(processorId, "terminated")
}
