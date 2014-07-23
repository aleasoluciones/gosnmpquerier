package snmpquery

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
	Oid         string
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
