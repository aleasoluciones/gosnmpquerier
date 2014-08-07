// Copyright 2014 The GoSNMPQuerier Authors. All rights reserved.  Use of this
// source code is governed by a MIT-style license that can be found in the
// LICENSE file.

package gosnmpquerier

import (
	"time"

	"github.com/eferro/gosnmp"
)

func walk(destination, community, oid string, timeout time.Duration, retries int) ([]gosnmp.SnmpPDU, error) {
	conn := snmpConnection(destination, community, timeout, retries)
	err := conn.Connect()
	if err != nil {
		return nil, err
	}
	defer conn.Conn.Close()
	output := make(chan gosnmp.SnmpPDU)
	errChannel := make(chan error, 1)
	go doWalk(conn, oid, output, errChannel)

	result := []gosnmp.SnmpPDU{}
	for pdu := range output {
		result = append(result, pdu)
	}
	if len(errChannel) != 0 {
		err := <-errChannel
		return nil, err
	}
	return result, nil
}

func doWalk(conn gosnmp.GoSNMP, oid string, output chan gosnmp.SnmpPDU, errChannel chan error) {
	processPDU := func(pdu gosnmp.SnmpPDU) error {
		output <- pdu
		return nil
	}
	err := conn.BulkWalk(oid, processPDU)
	if err != nil {
		errChannel <- err
	}
	close(output)
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

func get(destination, community string, oids []string, timeout time.Duration, retries int) ([]gosnmp.SnmpPDU, error) {
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

func getnext(destination, community string, oids []string, timeout time.Duration, retries int) ([]gosnmp.SnmpPDU, error) {
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
