package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type PageInfo struct {
	URL  string
	Size int   // количество байт в теле ответа
	Err  error // nil, если удалось получить
}

type res struct {
	index int
	PageInfo
}

type req struct {
	index   int
	url     string
	timeout time.Duration
}

func FetchPages(
	urls []string,
	concurrency int,
	timeout time.Duration,
) []PageInfo {
	in := make(chan req)
	out := make(chan res)
	var wg sync.WaitGroup

	go func() {
		defer close(in)
		for i := 0; i < len(urls); i++ {
			r := req{index: i, url: urls[i], timeout: timeout}
			in <- r
		}
	}()

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go worker(i, in, out, &wg)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	got := make([]res, 0)
	for o := range out {
		got = append(got, o)
	}
	ans := make([]PageInfo, len(got))
	for _, g := range got {
		ans[g.index] = g.PageInfo
	}
	return ans
}

func worker(id int, in <-chan req, out chan res, wg *sync.WaitGroup) {
	defer wg.Done()
	for t := range in {
		client := http.Client{
			Timeout: t.timeout,
		}
		resp, err := client.Get(t.url)
		if err != nil {
			pageInfo := PageInfo{URL: t.url, Size: 0, Err: err}
			res := res{index: t.index, PageInfo: pageInfo}
			out <- res
			continue
		}
		defer resp.Body.Close()
		body, read_err := io.ReadAll(resp.Body)
		if read_err != nil {
			fmt.Println("Error reading body:", read_err)
			return
		}
		pageInfo := PageInfo{URL: t.url, Size: len(body), Err: nil}
		res := res{index: t.index, PageInfo: pageInfo}
		out <- res
	}
}

