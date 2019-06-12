package pagination

import "fmt"

func ExampleMaxPageSizeFrom() {
	fmt.Println(maxPageSizeFrom(nil))
	fmt.Println(maxPageSizeFrom([]Option{{}}))
	fmt.Println(maxPageSizeFrom([]Option{{MaxPageSize: 100}}))

	// Output:
	// 1000
	// 1000
	// 100
}

func ExampleDefaultPageSizeFrom() {
	fmt.Println(defaultPageSizeFrom(nil))
	fmt.Println(defaultPageSizeFrom([]Option{{}}))
	fmt.Println(defaultPageSizeFrom([]Option{{DefaultPageSize: 100}}))

	// Output:
	// 10
	// 10
	// 100
}
