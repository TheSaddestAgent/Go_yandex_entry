package main

import (
	"sync"
)

type task struct {
	index int
	value int
}

func ParallelSquares(nums []int, workers int) []int {

	in := make(chan task)
	out := make(chan task)
	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go worker(i, in, out, &wg)
	}

	go func() {
		for i, num := range nums {
			t := task{index: i, value: num}
			in <- t
		}
		close(in)
	}()

	go func() {
		wg.Wait()
		close(out)
	}()

	got := make([]task, 0)
	for o := range out {
		got = append(got, o)
	}
	res := make([]int, len(got))
	for _, g := range got {
		res[g.index] = g.value
	}
	return res
}

func worker(id int, in <-chan task, out chan<- task, wg *sync.WaitGroup) {
	defer wg.Done()
	for t := range in {
		res_v := t.value * t.value
		res_i := t.index
		res := task{index: res_i, value: res_v}
		out <- res
	}
}

