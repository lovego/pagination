package pagination

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/lovego/errs"
)

type Pagination struct {
	TotalSize   int64 `json:"totalSize"`
	TotalPage   int64 `json:"totalPage"`
	CurrentPage int64 `json:"-"`
	PageSize    int64 `json:"-"`
}
type Querier interface {
	Query(data interface{}, sql string, args ...interface{}) error
}

func New(page, size string, maxPageSize int64) *Pagination {
	currentPage, _ := strconv.ParseInt(page, 10, 64)
	pageSize, _ := strconv.ParseInt(size, 10, 64)

	return NewFromInt64(currentPage, pageSize, maxPageSize)
}

func NewFromQuery(query url.Values, maxPageSize int64) *Pagination {
	return New(query.Get("page"), query.Get("pageSize"), maxPageSize)
}

func NewFromInt64(page, size, maxPageSize int64) *Pagination {
	if page <= 0 {
		page = 1
	}
	if maxPageSize <= 0 {
		maxPageSize = 10
	}
	if size <= 0 || size > maxPageSize {
		size = maxPageSize
	}
	return &Pagination{CurrentPage: page, PageSize: size}
}

func (p *Pagination) SQL() string {
	return fmt.Sprintf("LIMIT %d OFFSET %d", p.PageSize, (p.CurrentPage-1)*p.PageSize)
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

// we implemented sql.Scanner, so just use rows.Scan(pagination).
func (p *Pagination) Scan(src interface{}) error {
	switch totalSize := src.(type) {
	case int64:
		p.TotalSize = totalSize
		p.CalcTotalPage()
		return nil
	default:
		return fmt.Errorf("bsql: cannot assign %T(%v) to int64", src, src)
	}
}

func (p *Pagination) SetTotalSize(totalSize int) {
	p.TotalSize = int64(totalSize)
	p.CalcTotalPage()
}

func (p *Pagination) CalcTotalPage() {
	var totalPage = p.TotalSize / p.PageSize
	if p.TotalSize%p.PageSize > 0 {
		totalPage += 1
	}
	p.TotalPage = totalPage
}
