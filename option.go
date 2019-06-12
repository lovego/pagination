package pagination

var DefaultPageSize, MaxPageSize int64 = 10, 1000

type Option struct {
	DefaultPageSize int64
	MaxPageSize     int64
}

func maxPageSizeFrom(options []Option) int64 {
	if len(options) > 0 && options[0].MaxPageSize > 0 {
		return options[0].MaxPageSize
	}
	return MaxPageSize
}

func defaultPageSizeFrom(options []Option) int64 {
	if len(options) > 0 && options[0].DefaultPageSize > 0 {
		return options[0].DefaultPageSize
	}
	return DefaultPageSize
}
