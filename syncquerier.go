// Copyright 2014 The GoSNMPQuerier Authors. All rights reserved.  Use of this
// source code is governed by a MIT-style license that can be found in the
// LICENSE file.

package gosnmpquerier

import (
	"errors"
	"time"

	"github.com/soniah/gosnmp"
)

const (
	QUERIER_TIMEOUT = 1 * time.Second
)

type SyncQuerier interface {
	ExecuteQuery(query Query) Query
	Get(destination, community string, oids []string, timeout time.Duration, retries int) ([]gosnmp.SnmpPDU, error)
	GetNext(destination, community string, oids []string, timeout time.Duration, retries int) ([]gosnmp.SnmpPDU, error)
	Walk(destination, community, oid string, timeout time.Duration, retries int) ([]gosnmp.SnmpPDU, error)
}

type syncQuerier struct {
	Input        chan QueryWithOutputChannel
	asyncQuerier *AsyncQuerier
}

func NewSyncQuerier(contention, numErrors int, resetTime time.Duration) *syncQuerier {
	querier := syncQuerier{
		Input:        make(chan QueryWithOutputChannel),
		asyncQuerier: NewAsyncQuerier(contention, numErrors, resetTime),
	}
	go querier.processAndDispatchQueries()
	return &querier
}

func (querier *syncQuerier) ExecuteQuery(query Query) Query {
	output := make(chan Query)

	timeoutTimer := time.NewTimer(QUERIER_TIMEOUT)
	defer timeoutTimer.Stop()

	select {
	case querier.Input <- QueryWithOutputChannel{query, output}:

	// Same as time.After() but prevents memory leaks by manually stopping the timer
	case <-timeoutTimer.C:
		query.Error = errors.New("Global queries queue full")
		return query
	}

	return <-output
}

func (querier *syncQuerier) Get(destination, community string, oids []string, timeout time.Duration, retries int) ([]gosnmp.SnmpPDU, error) {
	return querier.executeCommand(GET, destination, community, oids, timeout, retries)
}

func (querier *syncQuerier) GetNext(destination, community string, oids []string, timeout time.Duration, retries int) ([]gosnmp.SnmpPDU, error) {
	return querier.executeCommand(GETNEXT, destination, community, oids, timeout, retries)
}

func (querier *syncQuerier) Walk(destination, community, oid string, timeout time.Duration, retries int) ([]gosnmp.SnmpPDU, error) {
	return querier.executeCommand(WALK, destination, community, []string{oid}, timeout, retries)
}

func (querier *syncQuerier) executeCommand(command OpSnmp, destination, community string, oids []string, timeout time.Duration, retries int) ([]gosnmp.SnmpPDU, error) {
	query := querier.makeQuery(command, destination, community, oids, timeout, retries)
	processedQuery := querier.ExecuteQuery(query)
	return processedQuery.Response, processedQuery.Error
}

func (querier *syncQuerier) makeQuery(command OpSnmp, destination, community string, oids []string, timeout time.Duration, retries int) Query {
	return Query{
		Cmd:         command,
		Community:   community,
		Oids:        oids,
		Timeout:     timeout,
		Retries:     retries,
		Destination: destination,
	}
}

func (querier *syncQuerier) processAndDispatchQueries() {
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
