package mocks

import "github.com/aleasoluciones/gosnmpquerier"
import "github.com/stretchr/testify/mock"

import "time"

import "github.com/soniah/gosnmp"

type SyncQuerier struct {
	mock.Mock
}

func (m *SyncQuerier) ExecuteQuery(query gosnmpquerier.Query) gosnmpquerier.Query {
	ret := m.Called(query)

	r0 := ret.Get(0).(gosnmpquerier.Query)

	return r0
}
func (m *SyncQuerier) Get(destination, community string, oids []string, timeout time.Duration, retries int) ([]gosnmp.SnmpPDU, error) {
	ret := m.Called(destination, community, oids, timeout, retries)

	r0 := ret.Get(0).([]gosnmp.SnmpPDU)
	r1 := ret.Error(1)

	return r0, r1
}
func (m *SyncQuerier) GetNext(destination, community string, oids []string, timeout time.Duration, retries int) ([]gosnmp.SnmpPDU, error) {
	ret := m.Called(destination, community, oids, timeout, retries)

	r0 := ret.Get(0).([]gosnmp.SnmpPDU)
	r1 := ret.Error(1)

	return r0, r1
}
func (m *SyncQuerier) Walk(destination, community, oid string, timeout time.Duration, retries int) ([]gosnmp.SnmpPDU, error) {
	ret := m.Called(destination, timeout, retries)

	r0 := ret.Get(0).([]gosnmp.SnmpPDU)
	r1 := ret.Error(1)

	return r0, r1
}
