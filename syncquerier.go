// Copyright 2014 The GoSNMPQuerier Authors. All rights reserved.  Use of this
// source code is governed by a MIT-style license that can be found in the
// LICENSE file.

package gosnmpquerier

import (
	"time"

	"github.com/eferro/gosnmp"
)

type SyncQuerier struct {
	Input        chan QueryWithOutputChannel
	asyncQuerier *AsyncQuerier
}

func NewSyncQuerier(contention int) *SyncQuerier {
	querier := SyncQuerier{
		Input:        make(chan QueryWithOutputChannel),
		asyncQuerier: NewAsyncQuerier(contention),
	}
	go querier.processAndDispatchQueries()
	return &querier
}

func (querier *SyncQuerier) ExecuteQuery(query Query) Query {
	output := make(chan Query)
	querier.Input <- QueryWithOutputChannel{query, output}
	processedQuery := <-output
	return processedQuery
}

func (querier *SyncQuerier) Get(destination, community string, oids []string, timeout time.Duration, retries int) ([]gosnmp.SnmpPDU, error) {
	query := Query{
		Cmd:         GET,
		Community:   community,
		Oids:        oids,
		Timeout:     timeout,
		Retries:     retries,
		Destination: destination,
	}

	processedQuery := querier.ExecuteQuery(query)
	return processedQuery.Response, processedQuery.Error
}

func (querier *SyncQuerier) GetNext(destination, community string, oids []string, timeout time.Duration, retries int) ([]gosnmp.SnmpPDU, error) {
	query := Query{
		Cmd:         GETNEXT,
		Community:   community,
		Oids:        oids,
		Timeout:     timeout,
		Retries:     retries,
		Destination: destination,
	}

	processedQuery := querier.ExecuteQuery(query)
	return processedQuery.Response, processedQuery.Error
}

func (querier *SyncQuerier) Walk(destination, community, oid string, timeout time.Duration, retries int) ([]gosnmp.SnmpPDU, error) {
	query := Query{
		Cmd:         WALK,
		Community:   community,
		Oids:        []string{oid},
		Timeout:     timeout,
		Retries:     retries,
		Destination: destination,
	}

	processedQuery := querier.ExecuteQuery(query)
	return processedQuery.Response, processedQuery.Error
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
