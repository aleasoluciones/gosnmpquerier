// Copyright 2014 The GoSNMPQuerier Authors. All rights reserved.  Use of this
// source code is governed by a MIT-style license that can be found in the
// LICENSE file.

package gosnmpquerier

import (
	"time"

	"github.com/soniah/gosnmp"
)

type SnmpClient struct {
}

func newSnmpClient() *SnmpClient {
	return &SnmpClient{}
}

func (snmpClient *SnmpClient) get(destination, community string, oids []string, timeout time.Duration, retries int) ([]gosnmp.SnmpPDU, error) {

	conn := snmpConnection(destination, community, timeout, retries)
	err := conn.Connect()
	if err != nil {
		return nil, err
	}
	defer conn.Conn.Close()

	result, err := conn.Get(oids)
	if err != nil {
		return nil, err
	}

	pdus := []gosnmp.SnmpPDU{}
	for _, pdu := range result.Variables {
		pdus = append(pdus, pdu)
	}
	return pdus, nil

}

func (snmpClient *SnmpClient) walk(destination, community, oid string, timeout time.Duration, retries int) ([]gosnmp.SnmpPDU, error) {
	conn := snmpConnection(destination, community, timeout, retries)
	err := conn.Connect()
	if err != nil {
		return nil, err
	}
	defer conn.Conn.Close()

	return conn.BulkWalkAll(oid)
}

func snmpConnection(destination, community string, timeout time.Duration, retries int) gosnmp.GoSNMP {
	return gosnmp.GoSNMP{
		Target:    destination,
		Port:      161,
		Community: community,
		Version:   gosnmp.Version2c,
		Timeout:   timeout,
		Retries:   retries,
	}
}

func (snmpClient *SnmpClient) getnext(destination, community string, oids []string, timeout time.Duration, retries int) ([]gosnmp.SnmpPDU, error) {
	conn := snmpConnection(destination, community, timeout, retries)
	err := conn.Connect()
	if err != nil {
		return nil, err
	}
	defer conn.Conn.Close()

	result, err := conn.GetNext(oids)
	if err != nil {
		return nil, err
	}

	pdus := []gosnmp.SnmpPDU{}
	for _, pdu := range result.Variables {
		pdus = append(pdus, pdu)
	}
	return pdus, nil
}
