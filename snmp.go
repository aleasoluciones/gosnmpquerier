// Copyright 2014 The GoSNMPQuerier Authors. All rights reserved.  Use of this
// source code is governed by a MIT-style license that can be found in the
// LICENSE file.

package gosnmpquerier

import (
	"time"

	"github.com/soniah/gosnmp"
)

type SnmpClient interface {
	get(destination, community string, oids []string, timeout time.Duration, retries int) ([]gosnmp.SnmpPDU, error)
	walk(destination, community, oid string, timeout time.Duration, retries int) ([]gosnmp.SnmpPDU, error)
	getnext(destination, community string, oids []string, timeout time.Duration, retries int) ([]gosnmp.SnmpPDU, error)
}

type GoSnmpClient struct{}

func newSnmpClient() *GoSnmpClient {
	return &GoSnmpClient{}
}

func (snmpClient *GoSnmpClient) get(destination, community string, oids []string, timeout time.Duration, retries int) ([]gosnmp.SnmpPDU, error) {
	conn := snmpConnection(destination, community, timeout, retries)
	if err := conn.Connect(); err != nil {
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

func (snmpClient *GoSnmpClient) walk(destination, community, oid string, timeout time.Duration, retries int) ([]gosnmp.SnmpPDU, error) {
	conn := snmpConnection(destination, community, timeout, retries)
	if err := conn.Connect(); err != nil {
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

func (snmpClient *GoSnmpClient) getnext(destination, community string, oids []string, timeout time.Duration, retries int) ([]gosnmp.SnmpPDU, error) {
	conn := snmpConnection(destination, community, timeout, retries)
	if err := conn.Connect(); err != nil {
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
