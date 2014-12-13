# GoSNMPQuerier

[![Build Status](https://travis-ci.org/aleasoluciones/gosnmpquerier.svg?branch=master)](https://travis-ci.org/aleasoluciones/gosnmpquerier)
[![Coverage Status](https://img.shields.io/coveralls/aleasoluciones/gosnmpquerier.svg)](https://coveralls.io/r/aleasoluciones/gosnmpquerier?branch=master)

Scalable SNMP querier library

## Features
 * Asynchronous and Synchronous scalable snmpquerier
 * Support for Walk/Get/GetNext snmp queries
 * Maximum number of concurrent snmp queries for each device/host (Contention).
 * Circuit Breaker pattern for each device/host connection
 * Back preasure control for each device/host incomming queries
 * Back preasure control for global incomming queries

##  Unimplemented features / TODO
 * Contention level configured for each device/host
 * Back preasure configuration
 * Set snmp command
 

## Installation

```
$ go get github.com/aleasoluciones/gosnmpquerier
```

Add it to your code:

```go
import "github.com/aleasoluciones/gosnmpquerier"
```

## Testing

```
$ go test -v
```

##License
(The MIT License)

Copyright (c) 2014 Alea Soluciones SLL

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the 'Software'), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED 'AS IS', WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

