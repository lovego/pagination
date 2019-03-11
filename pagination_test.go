package pagination

import "fmt"

func ExamplePagination() {
	pagination := New("3", "20", 50)
	fmt.Println(pagination.SQL())
	// we implemented sql.Scanner
	pagination.Scan(int64(82))
	fmt.Printf("%+v\n", pagination)

	// Output:
	// LIMIT 20 OFFSET 40
	// &{TotalSize:82 TotalPage:5 currentPage:3 pageSize:20}
}

func ExamplePagination_invalidParams() {
	pagination := New("", "100", 50)
	fmt.Println(pagination.SQL())
	pagination.Scan(int64(82))
	fmt.Printf("%+v\n", pagination)

	// Output:
	// LIMIT 50 OFFSET 0
	// &{TotalSize:82 TotalPage:2 currentPage:1 pageSize:50}
}