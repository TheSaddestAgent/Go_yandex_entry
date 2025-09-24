package main

import (
	"sync"
)

type Job func() error

type task struct {
	index int
	value Job
}
type result struct {
	index int
	res   error
}

func RunAllErrors(jobs []Job) []error {
	in := make(chan task)
	out := make(chan result)
	var wg sync.WaitGroup
	workers := len(jobs)
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go worker(i, in, out, &wg)
	}

	go func() {
		for i, j := range jobs {
			t := task{index: i, value: j}
			in <- t
		}
		close(in)
	}()

	go func() {
		wg.Wait()
		close(out)
	}()

	got := make([]result, 0)
	for o := range out {
		got = append(got, o)
	}
	res := make([]error, len(got))
	for _, g := range got {
		res[g.index] = g.res
	}
	return res
}

func worker(id int, in <-chan task, out chan<- result, wg *sync.WaitGroup) {
	defer wg.Done()
	for t := range in {
		var res_v error
		func() {
			defer func() {
				if r := recover(); r != nil {
				}
			}()

			res_v = t.value()
		}()
		res_i := t.index
		res := result{index: res_i, res: res_v}
		out <- res
	}
}
