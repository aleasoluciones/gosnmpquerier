package gosnmpquerier

import (
	"testing"
	"time"

	"github.com/gosnmp/gosnmp"
	"github.com/stretchr/testify/assert"
)

type FakeSnmpClient struct{}

func (snmpClient *FakeSnmpClient) get(destination, community string, oids []string, timeout time.Duration, retries int) ([]gosnmp.SnmpPDU, error) {
	return makeSnmpPDU()
}

func (snmpClient *FakeSnmpClient) walk(destination, community, oid string, timeout time.Duration, retries int) ([]gosnmp.SnmpPDU, error) {
	return makeSnmpPDU()
}

func (snmpClient *FakeSnmpClient) getnext(destination, community string, oids []string, timeout time.Duration, retries int) ([]gosnmp.SnmpPDU, error) {
	return makeSnmpPDU()
}

func makeSnmpPDU() ([]gosnmp.SnmpPDU, error) {
	return []gosnmp.SnmpPDU{{Name: "foo", Type: 1, Value: 1}}, nil
}

func newSyncQuerier() *syncQuerier {
	querier := NewSyncQuerier(1, 3, 3*time.Second)
	querier.asyncQuerier.snmpClient = &FakeSnmpClient{}
	return querier
}

func expectedSnmpResult() []gosnmp.SnmpPDU {
	return []gosnmp.SnmpPDU{{Name: "foo", Type: 0x1, Value: 1}}
}

func TestGetReturnsSnmpGetResult(t *testing.T) {
	querier := newSyncQuerier()
	result, _ := querier.Get("192.168.5.15", "alea2", []string{"1.3.6.1.2.1.1.1.0"}, 1*time.Second, 1)
	assert.Equal(t, result, expectedSnmpResult())
}

func TestGetNextReturnsSnmpGetNextResult(t *testing.T) {
	querier := newSyncQuerier()
	result, _ := querier.GetNext("192.168.5.15", "alea2", []string{"1.3.6.1.2.1.1.1.0"}, 1*time.Second, 1)
	assert.Equal(t, result, expectedSnmpResult())
}

func TestWalkReturnsSnmpWalkResult(t *testing.T) {
	querier := newSyncQuerier()
	result, _ := querier.Walk("192.168.5.15", "alea2", "1.3.6.1.2.1.1", 1*time.Second, 1)
	assert.Equal(t, result, expectedSnmpResult())
}
