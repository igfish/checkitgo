package checkit

type Iterator[T any] interface {
	Next() (T, bool)
}
