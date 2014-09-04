package gosnmpquerier

import (
	"fmt"
	"testing"
	"time"
	//"github.com/stretchr/testify/assert"
)

// SnmpClient
// SnmpClient.Get(destination, community, oids, duration, retries)

func TestGetReturnsSnmpGetResult(t *testing.T) {
	querier := NewSyncQuerier(1)
	result, err := querier.Get("192.168.5.15", "alea2", []string{"1.3.6.1.2.1.1.1.0"}, 1*time.Second, 1)
	// mockear el snmpwrapper que no tenemos
	fmt.Println(result, err)
}

// puedo hacer un Get y lo que me devuelve el wrapper snmp es lo que se me devuelve
// puedo hacer un Getnext y lo que me devuelve el wrapper snmp es lo que se me devuelve
// puedo hacer un walk y lo que me devuelve el wrapper snmp es lo que se me devuelve

// gosnmpquerier
// API asincrono
// API sincrono
// contentcion por destino....
// circuit breaker por destino

// QSnmp -> destino (por destino maximo concurrente)
