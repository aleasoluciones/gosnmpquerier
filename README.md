# gosnmpquerier

[![CI](https://github.com/aleasoluciones/gosnmpquerier/actions/workflows/ci.yml/badge.svg?branch=master)](https://github.com/aleasoluciones/gosnmpquerier/actions/workflows/ci.yml)
[![Coverage Status](https://img.shields.io/coveralls/aleasoluciones/gosnmpquerier.svg)](https://coveralls.io/r/aleasoluciones/gosnmpquerier?branch=master)
[![GoDoc](https://godoc.org/github.com/aleasoluciones/gosnmpquerier?status.png)](http://godoc.org/github.com/aleasoluciones/gosnmpquerier)
[![License](https://img.shields.io/github/license/aleasoluciones/gosnmpquerier)](https://github.com/aleasoluciones/gosnmpquerier/blob/master/LICENSE)

Scalable SNMP querier library

## Features

- Asynchronous and Synchronous scalable snmpquerier
- Support for Walk/Get/GetNext snmp queries
- Maximum number of concurrent snmp queries for each device/host (Contention).
- Circuit Breaker pattern for each device/host connection
- Back preasure control for each device/host incomming queries
- Back preasure control for global incomming queries

##  Unimplemented features / TODO

- Contention level configured for each device/host
- Back preasure configuration
- Set snmp command

## Testing

```
$ make test
```
