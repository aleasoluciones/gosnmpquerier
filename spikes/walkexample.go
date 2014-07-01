// Copyright 2014 Chris Dance (codedance). All rights reserved.  Use of this
// source code is governed by a BSD-style license that can be found in
// the LICENSE file.

// This program demonstrates BulkWalk.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/soniah/gosnmp"
)

func usage() {
	fmt.Println("Usage:")
	fmt.Printf("   %s community host [oid]\n", filepath.Base(os.Args[0]))
	fmt.Println("     community - community")
	fmt.Println("     host - the host to walk/scan")
	fmt.Println("     oid  - the MIB/Oid defining a subtree of values")
	os.Exit(1)
}

func main() {
	if len(os.Args) < 4 {
		usage()
	}
	community := os.Args[1]
	target := os.Args[2]
	var oid string
	if len(os.Args) > 3 {
		oid = os.Args[3]
	}

	gosnmp.Default.Community = community
	gosnmp.Default.Target = target
	gosnmp.Default.Timeout = time.Duration(10 * time.Second) // Timeout better suited to walking
	err := gosnmp.Default.Connect()
	if err != nil {
		fmt.Printf("Connect err: %v\n", err)
		os.Exit(1)
	}
	defer gosnmp.Default.Conn.Close()

	err = gosnmp.Default.BulkWalk(oid, printValue)
	if err != nil {
		fmt.Printf("Walk Error: %v\n", err)
		os.Exit(1)
	}
}

func printValue(pdu gosnmp.SnmpPDU) error {
	fmt.Printf("%s = ", pdu.Name)

	switch pdu.Type {
	case gosnmp.OctetString:
		fmt.Printf("STRING: %s\n", pdu.Value.(string))
	case gosnmp.Counter64:
		fmt.Printf("COUNTER64 %d: %d\n", pdu.Type, gosnmp.ToBigInt(pdu.Value))
	default:
		fmt.Printf("TYPE %d: %d\n", pdu.Type, gosnmp.ToBigInt(pdu.Value))
	}
	return nil
}
