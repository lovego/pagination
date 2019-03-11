package pagination

import (
	"fmt"
	"strconv"
)

type Pagination struct {
	TotalSize   int64 `json:"totalSize"`
	TotalPage   int64 `json:"totalPage"`
	currentPage int64
	pageSize    int64 `json:"pageSize"`
}

func New(page, size string, maxPageSize int64) *Pagination {
	currentPage, _ := strconv.ParseInt(page, 10, 64)
	if currentPage <= 0 {
		currentPage = 1
	}

	pageSize, _ := strconv.ParseInt(size, 10, 64)
	if pageSize <= 0 || pageSize > maxPageSize {
		pageSize = maxPageSize
	}

	return &Pagination{currentPage: currentPage, pageSize: pageSize}
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
