# pagination
pagination util for golang.

[![Build Status](https://github.com/lovego/pagination/actions/workflows/go.yml/badge.svg)](https://github.com/lovego/pagination/actions/workflows/go.yml)
[![Coverage Status](https://coveralls.io/repos/github/lovego/pagination/badge.svg?branch=master&1)](https://coveralls.io/github/lovego/pagination)
[![Go Report Card](https://goreportcard.com/badge/github.com/lovego/pagination)](https://goreportcard.com/report/github.com/lovego/pagination)
[![Documentation](https://pkg.go.dev/badge/github.com/lovego/pagination)](https://pkg.go.dev/github.com/lovego/pagination@v0.0.1)

## Install
`$ go get github.com/lovego/pagination`

## Usage
```go
package pagination_test

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
	result := Result{Pagination: pagination.New(page, size)}

	var mainSql = fmt.Sprintf("FROM users WHERE type = %d", 1)
	if err := db.Query(&result.Rows, fmt.Sprintf(
		"SELECT id, name %s %s", mainSql, result.Pagination.SQL(),
	)); err != nil {
		return nil, errs.Trace(err)
	}
	if err := result.SetupTotalSize(len(result.Rows), db, "SELECT count(*) "+mainSql); err != nil {
		return nil, errs.Trace(err)
	}

	return &result, nil
}

func List2(page, size string) (*Result, error) {
	result := Result{
		Pagination: pagination.New(page, size, pagination.Option{DefaultPageSize: 10, MaxPageSize: 100}),
	}

	var mainSql = fmt.Sprintf("FROM users WHERE type = %d", 1)
	if err := db.Query(&result.Rows, fmt.Sprintf(
		"SELECT id, name %s %s", mainSql, result.Pagination.SQL(),
	)); err != nil {
		return nil, errs.Trace(err)
	}

	if result.CurrentPage == 1 && len(result.Rows) < int(result.PageSize) {
		result.SetTotalSize(len(result.Rows))
	} else {
		if err := db.Query(result, "SELECT count(*) "+mainSql); err != nil {
			return nil, errs.Trace(err)
		}
	}

	return &result, nil
}
```

## Documentation
[https://godoc.org/github.com/lovego/pagination](https://godoc.org/github.com/lovego/pagination)
