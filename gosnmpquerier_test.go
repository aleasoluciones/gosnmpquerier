package gosnmpquerier

import (
	"testing"
	"time"

	"github.com/soniah/gosnmp"
	"github.com/stretchr/testify/assert"
)

func TestGetReturnsSnmpGetResult(t *testing.T) {
	querier := NewSyncQuerier(1)
	result, _ := querier.Get("192.168.5.15", "alea2", []string{"1.3.6.1.2.1.1.1.0"}, 1*time.Second, 1)
	assert.Equal(t, result, []gosnmp.SnmpPDU{gosnmp.SnmpPDU{Name: ".1.3.6.1.2.1.1.1.0", Type: 0x4, Value: []uint8{0x48, 0x75, 0x61, 0x77, 0x65, 0x69, 0x20, 0x49, 0x6e, 0x74, 0x65, 0x67, 0x72, 0x61, 0x74, 0x65, 0x64, 0x20, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x20, 0x53, 0x6f, 0x66, 0x74, 0x77, 0x61, 0x72, 0x65}}})
}

func TestGetNextReturnsSnmpGetNextResult(t *testing.T) {
	querier := NewSyncQuerier(1)
	result, _ := querier.GetNext("192.168.5.15", "alea2", []string{"1.3.6.1.2.1.1.1.0"}, 1*time.Second, 1)
	assert.Equal(t, result, []gosnmp.SnmpPDU{gosnmp.SnmpPDU{Name: ".1.3.6.1.2.1.1.2.0", Type: 0x6, Value: ".1.3.6.1.4.1.2011.2.123"}})
}

// puedo hacer un walk y lo que me devuelve el wrapper snmp es lo que se me devuelve

// gosnmpquerier
// API asincrono
// API sincrono
// contentcion por destino....
// circuit breaker por destino

// QSnmp -> destino (por destino maximo concurrente)
