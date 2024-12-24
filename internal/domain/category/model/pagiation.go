package model

type Pagination[T any] struct {
	Offset int64
	Limit  int64
	Total  int64
	Data   []T
}
