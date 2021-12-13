package pagination

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/lovego/errs"
)

type Paginator struct {
	Page     int64 `json:"page" c:"页码"`
	PageSize int64 `json:"pageSize" c:"页宽"`
}

func (p Paginator) Pagination(options ...Option) *Pagination {
	return NewFromInt64(p.Page, p.PageSize, options...)
}

type Pagination struct {
	TotalSize   int64 `json:"totalSize"`
	TotalPage   int64 `json:"totalPage"`
	CurrentPage int64 `json:"-"`
	PageSize    int64 `json:"-"`
}

type Querier interface {
	Query(data interface{}, sql string, args ...interface{}) error
	QueryCtx(ctx context.Context, opName string, data interface{}, sql string, args ...interface{}) error
}

// NewFromQuery returns a pagination from url.Values
func NewFromQuery(query url.Values, options ...Option) *Pagination {
	return New(query.Get("page"), query.Get("pageSize"), options...)
}

// New returns a Pagination from page, size in string type
func New(page, size string, options ...Option) *Pagination {
	currentPage, _ := strconv.ParseInt(page, 10, 64)
	pageSize, _ := strconv.ParseInt(size, 10, 64)

	return NewFromInt64(currentPage, pageSize, options...)
}

func NewFromInt64(page, size int64, options ...Option) *Pagination {
	if page <= 0 {
		page = 1
	}
	if size <= 0 || size > maxPageSizeFrom(options) {
		size = defaultPageSizeFrom(options)
	}
	return &Pagination{CurrentPage: page, PageSize: size}
}

func (p *Pagination) SQL() string {
	return fmt.Sprintf("LIMIT %d OFFSET %d", p.PageSize, p.Offset())
}

func (p *Pagination) Offset() int64 {
	return (p.CurrentPage - 1) * p.PageSize
}

func (p *Pagination) SetupTotalSize(
	firstPageSize int, querier Querier, sql string, args ...interface{},
) error {
	if p.CurrentPage == 1 && firstPageSize < int(p.PageSize) {
		p.SetTotalSize(firstPageSize)
		return nil
	}
	if err := querier.Query(p, sql, args...); err != nil {
		return errs.Trace(err)
	}
	return nil
}

func (p *Pagination) SetupTotalSizeCtx(
	ctx context.Context, opName string, firstPageSize int, querier Querier, sql string, args ...interface{},
) error {
	if p.CurrentPage == 1 && firstPageSize < int(p.PageSize) {
		p.SetTotalSize(firstPageSize)
		return nil
	}
	if err := querier.QueryCtx(ctx, opName, p, sql, args...); err != nil {
		return errs.Trace(err)
	}
	return nil
}

func (p *Pagination) SetupTotalSizeFunc(
	firstPageSize int, querier Querier, sqlFunc func() (string, error), args ...interface{},
) error {
	if p.CurrentPage == 1 && firstPageSize < int(p.PageSize) {
		p.SetTotalSize(firstPageSize)
		return nil
	}
	sql, err := sqlFunc()
	if err != nil {
		return err
	}
	if err := querier.Query(p, sql, args...); err != nil {
		return errs.Trace(err)
	}
	return nil
}

// Scan implemented sql.Scanner, so just use rows.Scan(pagination).
func (p *Pagination) Scan(src interface{}) error {
	switch totalSize := src.(type) {
	case int64:
		p.TotalSize = totalSize
		p.CalcTotalPage()
		return nil
	default:
		return fmt.Errorf("pagination: cannot assign %T(%v) to int64", src, src)
	}
}

func (p *Pagination) SetTotalSize(totalSize int) {
	p.TotalSize = int64(totalSize)
	p.CalcTotalPage()
}

func (p *Pagination) CalcTotalPage() {
	var totalPage = p.TotalSize / p.PageSize
	if p.TotalSize%p.PageSize > 0 {
		totalPage++
	}
	p.TotalPage = totalPage
}
