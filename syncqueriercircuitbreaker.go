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

type SyncQuerierCircuitBreaker struct {
	syncQuerier    *SyncQuerier
	circuitBreaker *circuitbreaker.Circuit
}

func NewSyncQuerierCircuitBreaker(contention, numErrors int, resetTime time.Duration) *SyncQuerierCircuitBreaker {
	return &SyncQuerierCircuitBreaker{
		syncQuerier:    NewSyncQuerier(contention),
		circuitBreaker: circuitbreaker.NewCircuit(numErrors, resetTime),
	}
}

func (querier *SyncQuerierCircuitBreaker) Get(destination, community string, oids []string, timeout time.Duration, retries int) ([]gosnmp.SnmpPDU, error) {
	if querier.circuitBreaker.IsOpen() {
		return nil, fmt.Errorf("Destination device unavailable %s", destination)
	}
	result, err := querier.syncQuerier.Get(destination, community, oids, timeout, retries)
	querier.reportCircuitStatus(err)
	return result, err
}

func (querier *SyncQuerierCircuitBreaker) reportCircuitStatus(err error) {
	if err == nil {
		querier.circuitBreaker.Ok()
	} else {
		querier.circuitBreaker.Error()
	}
}
