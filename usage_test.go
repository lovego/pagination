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
	result := Result{Pagination: pagination.New(page, size, 20)}

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
	result := Result{Pagination: pagination.New(page, size, 20)}

	var mainSql = fmt.Sprintf("FROM users WHERE type = %d", 1)
	if err := db.Query(&result.Rows, fmt.Sprintf(
		"SELECT id, name %s %s", mainSql, result.Pagination.SQL(),
	)); err != nil {
		return nil, errs.Trace(err)
	}

	if len(result.Rows) < int(result.PageSize) {
		result.SetTotalSize(len(result.Rows))
	} else {
		if err := db.Query(result, "SELECT count(*) "+mainSql); err != nil {
			return nil, errs.Trace(err)
		}
	}

	return &result, nil
}
