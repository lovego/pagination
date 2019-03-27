# pagination
pagination util for golang.

[![Build Status](https://travis-ci.org/lovego/pagination.svg?branch=master)](https://travis-ci.org/lovego/pagination)
[![Coverage Status](https://img.shields.io/coveralls/github/lovego/pagination/master.svg)](https://coveralls.io/github/lovego/pagination?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/lovego/pagination)](https://goreportcard.com/report/github.com/lovego/pagination)
[![GoDoc](https://godoc.org/github.com/lovego/pagination?status.svg)](https://godoc.org/github.com/lovego/pagination)

## Install
`$ go get github.com/lovego/pagination`

## Usage
```go
package example

import (
	"fmt"

	"github.com/lovego/bsql"
	"github.com/lovego/errs"
	"github.com/lovego/pagination"
)

var db *bsql.DB

type Result struct {
	Rows []Row `json:"rows"`
	*pagination.Pagination
}

type Row struct {
	Id   int64
	Name string
}

func List(page, size string) (*Result, error) {
	result := Result{Pagination: pagination.New(page, size, 20)}

	if err := db.Query(&result.Rows, fmt.Sprintf(
		"SELECT id, name FROM users %s", result.Pagination.SQL(),
	)); err != nil {
		return nil, errs.Trace(err)
	}

	if err := db.Query(result.Pagination, fmt.Sprintf("SELECT count(*) FROM users")); err != nil {
		return nil, errs.Trace(err)
	}

	return &result, nil
}
```

## Documentation
[https://godoc.org/github.com/lovego/pagination](https://godoc.org/github.com/lovego/pagination)
