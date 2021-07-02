package pagination

import (
	"context"
	"database/sql"
	"fmt"
)

func ExampleNew() {
	pagination := New("3", "20")
	fmt.Println(pagination.SQL())
	// we implemented sql.Scanner
	pagination.Scan(int64(82))
	fmt.Printf("%+v\n", pagination)

	// Output:
	// LIMIT 20 OFFSET 40
	// &{TotalSize:82 TotalPage:5 CurrentPage:3 PageSize:20}
}

func ExampleNew_2() {
	pagination := New("", "100", Option{DefaultPageSize: 10, MaxPageSize: 50})
	fmt.Println(pagination.SQL())
	pagination.Scan(int64(82))
	fmt.Printf("%+v\n", pagination)

	// Output:
	// LIMIT 10 OFFSET 0
	// &{TotalSize:82 TotalPage:9 CurrentPage:1 PageSize:10}
}

func ExampleNewFromQuery() {
	fmt.Printf("%+v\n", NewFromQuery(nil))
	// Output:
	// &{TotalSize:0 TotalPage:0 CurrentPage:1 PageSize:10}
}

type testQuerier struct{}

func (q testQuerier) Query(data interface{}, querySql string, args ...interface{}) error {
	if scanner, ok := data.(sql.Scanner); ok {
		scanner.Scan(int64(150))
	}
	return nil
}

func (q testQuerier) QueryCtx(ctx context.Context, opName string, data interface{}, querySql string, args ...interface{}) error {
	if scanner, ok := data.(sql.Scanner); ok {
		scanner.Scan(int64(150))
	}
	return nil
}

func ExamplePagination_SetupTotalSize() {
	p := NewFromInt64(1, 100)
	fmt.Println(p.SetupTotalSize(50, nil, "SELECT count(*)"))
	fmt.Printf("%+v\n", p)

	p = NewFromInt64(2, 100)
	fmt.Println(p.SetupTotalSize(50, testQuerier{}, "SELECT count(*)"))
	fmt.Printf("%+v\n", p)
	// Output:
	// <nil>
	// &{TotalSize:50 TotalPage:1 CurrentPage:1 PageSize:100}
	// <nil>
	// &{TotalSize:150 TotalPage:2 CurrentPage:2 PageSize:100}
}

func ExamplePagination_SetupTotalSizeCtx() {
	p := NewFromInt64(1, 100)
	fmt.Println(p.SetupTotalSizeCtx(context.Background(), `SetupTotalSizeCtx1`, 50, nil, "SELECT count(*)"))
	fmt.Printf("%+v\n", p)

	p = NewFromInt64(2, 100)
	fmt.Println(p.SetupTotalSizeCtx(context.Background(), `SetupTotalSizeCtx2`, 50, testQuerier{}, "SELECT count(*)"))
	fmt.Printf("%+v\n", p)
	// Output:
	// <nil>
	// &{TotalSize:50 TotalPage:1 CurrentPage:1 PageSize:100}
	// <nil>
	// &{TotalSize:150 TotalPage:2 CurrentPage:2 PageSize:100}
}

func ExamplePagination_SetupTotalSizeFunc() {
	p := NewFromInt64(1, 100)
	fmt.Println(p.SetupTotalSizeFunc(50, nil, func() (string, error) {
		return "", nil
	}))
	fmt.Printf("%+v\n", p)

	p = NewFromInt64(2, 100)
	fmt.Println(p.SetupTotalSizeFunc(50, testQuerier{}, func() (string, error) {
		return "", nil
	}))
	fmt.Printf("%+v\n", p)
	// Output:
	// <nil>
	// &{TotalSize:50 TotalPage:1 CurrentPage:1 PageSize:100}
	// <nil>
	// &{TotalSize:150 TotalPage:2 CurrentPage:2 PageSize:100}
}

func ExamplePagination_Scan() {
	p := NewFromInt64(1, 10)
	fmt.Println(p.Scan(10))
	// Output:
	// pagination: cannot assign int(10) to int64
}
