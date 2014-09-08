// Copyright 2014 The GoSNMPQuerier Authors. All rights reserved.  Use of this
// source code is governed by a MIT-style license that can be found in the
// LICENSE file.

package main

import (
	"fmt"
	"time"

	"github.com/aleasoluciones/gosnmpquerier"
)

const (
	CONTENTION = 4
)

func main() {

	querier := gosnmpquerier.NewSyncQuerier(CONTENTION, 3, 3*time.Second)

	for id := 0; id < 10; id++ {
		q, _ := gosnmpquerier.FromJson(`{"command":"walk", "destination":"localhost", "community":"public", "oids":["1.3.6.1"]}`)
		q.Id = id

		fmt.Println("Result:", querier.ExecuteQuery(*q))
	}
}
