// Copyright 2014 The GoSNMPQuerier Authors. All rights reserved.  Use of this
// source code is governed by a MIT-style license that can be found in the
// LICENSE file.

package gosnmpquerier

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/aleasoluciones/goaleasoluciones/circuitbreaker"
)

type AsyncQuerier struct {
	Input      chan Query
	Output     chan Query
	Contention int
	snmpClient SnmpClient
	numErrors  int
	resetTime  time.Duration
}

func NewAsyncQuerier(contention int, numErrors int, resetTime time.Duration) *AsyncQuerier {
	querier := AsyncQuerier{
		Input:      make(chan Query, 10),
		Output:     make(chan Query, 10),
		Contention: contention,
		snmpClient: newSnmpClient(),
		numErrors:  numErrors,
		resetTime:  resetTime,
	}
	go querier.process()
	return &querier
}

type destinationProcessor struct {
	querier        *AsyncQuerier
	input          chan Query
	output         chan Query
	done           chan bool
	circuitBreaker *circuitbreaker.Circuit
}

func (querier *AsyncQuerier) process() {
	log.Println("AsyncQuerier process begin")
	m := make(map[string]destinationProcessor)

	for query := range querier.Input {
		if _, exists := m[query.Destination]; exists == false {
			processorInfo := createProcessorInfo(querier, querier.Output)
			m[query.Destination] = processorInfo
			createProcessors(processorInfo, query.Destination)
		}
		m[query.Destination].input <- query
	}
	log.Println("AsyncQuerier process terminating")

	waitUntilProcessorEnd(m, querier.Contention)
	log.Println("closing output")
	close(querier.Output)
}
func createProcessorInfo(querier *AsyncQuerier, output chan Query) destinationProcessor {

	return destinationProcessor{
		querier:        querier,
		input:          make(chan Query, 10),
		output:         output,
		done:           make(chan bool, 1),
		circuitBreaker: circuitbreaker.NewCircuit(querier.numErrors, querier.resetTime),
	}
}

func createProcessors(processorInfo destinationProcessor, destination string) {
	for i := 0; i < processorInfo.querier.Contention; i++ {
		go processorInfo.processQueriesFromChannel(string(destination) + string("_") + strconv.Itoa(i))
	}
}

func waitUntilProcessorEnd(m map[string]destinationProcessor, contention int) {
	for destination, processorInfo := range m {
		log.Println("closing:", processorInfo.input)
		close(processorInfo.input)
		for i := 0; i < contention; i++ {
			<-processorInfo.done
		}
		delete(m, destination)
	}
}

func (processor *destinationProcessor) handleQuery(query *Query) {

	if processor.circuitBreaker.IsOpen() {
		query.Error = fmt.Errorf("destination device unavailable %s", query.Destination)
	} else {
		switch query.Cmd {
		case WALK:
			if len(query.Oids) == 1 {
				query.Response, query.Error = processor.querier.snmpClient.walk(query.Destination, query.Community, query.Oids[0], query.Timeout, query.Retries)
			} else {
				query.Error = fmt.Errorf("multi Oid Walk not supported")
			}
		case GET:
			query.Response, query.Error = processor.querier.snmpClient.get(query.Destination, query.Community, query.Oids, query.Timeout, query.Retries)
		case GETNEXT:
			query.Response, query.Error = processor.querier.snmpClient.getnext(query.Destination, query.Community, query.Oids, query.Timeout, query.Retries)
		}

		if query.Error != nil {
			processor.circuitBreaker.Error()
		} else {
			processor.circuitBreaker.Ok()
		}
	}
}

func (processor *destinationProcessor) processQueriesFromChannel(processorId string) {
	for query := range processor.input {
		processor.handleQuery(&query)
		processor.output <- query
	}
	processor.done <- true
	log.Println(processorId, "terminated")
}
