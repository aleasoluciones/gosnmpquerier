package gosnmpquerier

import (
	"testing"
	"time"

	"github.com/soniah/gosnmp"
	"github.com/stretchr/testify/assert"
)

type FakeSnmpClient struct{}

func (snmpClient *FakeSnmpClient) get(destination, community string, oids []string, timeout time.Duration, retries int) ([]gosnmp.SnmpPDU, error) {
	return []gosnmp.SnmpPDU{gosnmp.SnmpPDU{Name: "foo", Type: 1, Value: 1}}, nil
}

func (snmpClient *FakeSnmpClient) walk(destination, community, oid string, timeout time.Duration, retries int) ([]gosnmp.SnmpPDU, error) {
	return []gosnmp.SnmpPDU{gosnmp.SnmpPDU{Name: "foo", Type: 1, Value: 1}}, nil
}

func (snmpClient *FakeSnmpClient) getnext(destination, community string, oids []string, timeout time.Duration, retries int) ([]gosnmp.SnmpPDU, error) {
	return []gosnmp.SnmpPDU{gosnmp.SnmpPDU{Name: "foo", Type: 1, Value: 1}}, nil

}

func syncQuerier() *SyncQuerier {
	querier := NewSyncQuerier(1)
	querier.asyncQuerier.snmpClient = &FakeSnmpClient{}
	return querier
}

func TestGetReturnsSnmpGetResult(t *testing.T) {
	querier := syncQuerier()
	result, _ := querier.Get("192.168.5.15", "alea2", []string{"1.3.6.1.2.1.1.1.0"}, 1*time.Second, 1)
	assert.Equal(t, result, []gosnmp.SnmpPDU{gosnmp.SnmpPDU{Name: "foo", Type: 0x1, Value: 1}})
}

func TestGetNextReturnsSnmpGetNextResult(t *testing.T) {
	querier := syncQuerier()
	result, _ := querier.GetNext("192.168.5.15", "alea2", []string{"1.3.6.1.2.1.1.1.0"}, 1*time.Second, 1)
	assert.Equal(t, result, []gosnmp.SnmpPDU{gosnmp.SnmpPDU{Name: "foo", Type: 0x1, Value: 1}})
}

func TestWalkReturnsSnmpWalkResult(t *testing.T) {
	querier := syncQuerier()
	result, _ := querier.Walk("192.168.5.15", "alea2", "1.3.6.1.2.1.1", 1*time.Second, 1)
	assert.Equal(t, result, []gosnmp.SnmpPDU{gosnmp.SnmpPDU{Name: "foo", Type: 0x1, Value: 1}})
}
