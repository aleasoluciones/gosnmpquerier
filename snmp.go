// Copyright 2014 The GoSNMPQuerier Authors. All rights reserved.  Use of this
// source code is governed by a MIT-style license that can be found in the
// LICENSE file.

package gosnmpquerier

import (
	"log"
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
	t1 := time.Now()
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
	log.Println("get time", time.Since(t1), "data", destination, oids)
	return pdus, nil

}

func (snmpClient *GoSnmpClient) walk(destination, community, oid string, timeout time.Duration, retries int) ([]gosnmp.SnmpPDU, error) {
	t1 := time.Now()
	conn := snmpConnection(destination, community, timeout, retries)
	if err := conn.Connect(); err != nil {
		return nil, err
	}
	defer conn.Conn.Close()

	result, err := conn.BulkWalkAll(oid)
	log.Println("walk time", time.Since(t1), "data", destination, oid)
	return result, err
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
	t1 := time.Now()
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
	log.Println("getnext time", time.Since(t1), "data", destination, oids)
	return pdus, nil
}
