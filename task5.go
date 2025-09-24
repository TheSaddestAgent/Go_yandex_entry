package main

type Mapper[T any, R any] func(T) R

func MapChan[T any, R any](in <-chan T, fn Mapper[T, R]) <-chan R {
	out := make(chan R)
	go func() {
		defer close(out)
		for v := range in {
			out <- fn(v)
		}
	}()
	return out
}
