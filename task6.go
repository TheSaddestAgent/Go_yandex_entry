package main

import "sync"

func FanIn[T any](in1, in2 <-chan T) <-chan T {
	out := make(chan T)
	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		defer wg.Done()
		for v := range in1 {
			out <- v
		}
	}()

	go func() {
		defer wg.Done()
		for v := range in2 {
			out <- v
		}
	}()

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
