// Copyright 2014 The GoSNMPQuerier Authors. All rights reserved.  Use of this
// source code is governed by a MIT-style license that can be found in the
// LICENSE file.

package gosnmpquerier

import (
	"time"

	"github.com/eferro/gosnmp"
)

type OpSnmp int32

const (
	GET  = 0
	WALK = 1
)

type Query struct {
	Id          int
	Cmd         OpSnmp
	Community   string
	Oids        []string
	Timeout     time.Duration
	Retries     int
	Destination string
	Response    []gosnmp.SnmpPDU
	Error       error
}

type QueryWithOutputChannel struct {
	query           Query
	responseChannel chan Query
}
