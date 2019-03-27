package pagination

import (
	"fmt"
	"strconv"
)

type Pagination struct {
	TotalSize   int64 `json:"totalSize"`
	TotalPage   int64 `json:"totalPage"`
	currentPage int64
	pageSize    int64
}

func New(page, size string, maxPageSize int64) *Pagination {
	currentPage, _ := strconv.ParseInt(page, 10, 64)
	pageSize, _ := strconv.ParseInt(size, 10, 64)

	return NewFromInt64(currentPage, pageSize, maxPageSize)
}

func NewFromInt64(page, size, maxPageSize int64) *Pagination {
	if page <= 0 {
		page = 1
	}
	if size <= 0 || size > maxPageSize && maxPageSize > 0 {
		size = maxPageSize
	}
	return &Pagination{currentPage: page, pageSize: size}
}

func (p *Pagination) SQL() string {
	return fmt.Sprintf("LIMIT %d OFFSET %d", p.pageSize, (p.currentPage-1)*p.pageSize)
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

func (p *Pagination) CalcTotalPage() {
	var totalPage = p.TotalSize / p.pageSize
	if p.TotalSize%p.pageSize > 0 {
		totalPage += 1
	}
	p.TotalPage = totalPage
}
