// Copyright 2014 The GoSNMPQuerier Authors. All rights reserved.  Use of this
// source code is governed by a MIT-style license that can be found in the
// LICENSE file.

package gosnmpquerier

import (
	"fmt"
	"time"

	"github.com/aleasoluciones/gocircuitbreaker"
	"github.com/soniah/gosnmp"
)

type SyncQuerier struct {
	Input          chan QueryWithOutputChannel
	asyncQuerier   *AsyncQuerier
	circuitBreaker *circuitbreaker.Circuit
}

func NewSyncQuerier(contention, numErrors int, resetTime time.Duration) *SyncQuerier {
	querier := SyncQuerier{
		Input:          make(chan QueryWithOutputChannel),
		asyncQuerier:   NewAsyncQuerier(contention),
		circuitBreaker: circuitbreaker.NewCircuit(numErrors, resetTime),
	}
	go querier.processAndDispatchQueries()
	return &querier
}

func (querier *SyncQuerier) ExecuteQuery(query Query) Query {
	output := make(chan Query)
	querier.Input <- QueryWithOutputChannel{query, output}
	return <-output
}

func (querier *SyncQuerier) Get(destination, community string, oids []string, timeout time.Duration, retries int) ([]gosnmp.SnmpPDU, error) {
	return querier.executeCommand(GET, destination, community, oids, timeout, retries)
}

func (querier *SyncQuerier) GetNext(destination, community string, oids []string, timeout time.Duration, retries int) ([]gosnmp.SnmpPDU, error) {
	return querier.executeCommand(GETNEXT, destination, community, oids, timeout, retries)
}

func (querier *SyncQuerier) Walk(destination, community, oid string, timeout time.Duration, retries int) ([]gosnmp.SnmpPDU, error) {
	return querier.executeCommand(WALK, destination, community, []string{oid}, timeout, retries)
}

func (querier *SyncQuerier) executeCommand(command OpSnmp, destination, community string, oids []string, timeout time.Duration, retries int) ([]gosnmp.SnmpPDU, error) {
	if querier.circuitBreaker.IsOpen() {
		return nil, fmt.Errorf("Destination device unavailable %s", destination)
	}
	query := querier.makeQuery(command, destination, community, oids, timeout, retries)
	processedQuery := querier.ExecuteQuery(query)
	querier.reportCircuitStatus(processedQuery.Error)
	return processedQuery.Response, processedQuery.Error
}

func (querier *SyncQuerier) makeQuery(command OpSnmp, destination, community string, oids []string, timeout time.Duration, retries int) Query {
	return Query{
		Cmd:         command,
		Community:   community,
		Oids:        oids,
		Timeout:     timeout,
		Retries:     retries,
		Destination: destination,
	}
}

func (querier *SyncQuerier) reportCircuitStatus(err error) {
	if err == nil {
		querier.circuitBreaker.Ok()
	} else {
		querier.circuitBreaker.Error()
	}
}

func (querier *SyncQuerier) processAndDispatchQueries() {
	m := make(map[int]chan Query)
	i := 0
	for {
		select {
		case queryWithOutputChannel := <-querier.Input:
			queryWithOutputChannel.query.Id = i
			i += 1
			m[queryWithOutputChannel.query.Id] = queryWithOutputChannel.responseChannel
			querier.asyncQuerier.Input <- queryWithOutputChannel.query
		case processedQuery := <-querier.asyncQuerier.Output:
			m[processedQuery.Id] <- processedQuery
			delete(m, processedQuery.Id)
		}
	}
}
