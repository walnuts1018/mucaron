package result

type Result[T any] struct {
	Result T
	Error  error
}
