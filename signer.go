package main

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

func ExecutePipeline(jobs ...job) {
	wg := &sync.WaitGroup{}
	in := make(chan interface{}, 3)
	out := make(chan interface{}, 3)

	for i := 0; i < len(jobs); i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup, in, out chan interface{}, j job) {
			j(in, out)
			close(out)
			wg.Done()
		}(wg, in, out, jobs[i])
		in = out
		out = make(chan interface{}, 3)
	}
	//fmt.Scan()
	wg.Wait()
}

// SingleHash

func SingleHash(in, out chan interface{}) {
	m := &sync.Mutex{}
	wg := &sync.WaitGroup{}

	for i := range in {
		wg.Add(1)
		go singleHashWorker(wg, m, out, fmt.Sprintf("%v", i))
	}
	wg.Wait()
}

func singleHashWorker(wg *sync.WaitGroup, m *sync.Mutex, out chan interface{}, data string) {
	defer wg.Done()
	hash1 := make(chan string)
	hash2 := make(chan string)

	go func() {
		defer close(hash1)
		hash1 <- DataSignerCrc32(data)
	}()

	go func() {
		defer close(hash2)
		m.Lock()
		mb5 := DataSignerMd5(data)
		m.Unlock()
		hash2 <- DataSignerCrc32(mb5)
	}()

	out <- fmt.Sprintf("%s~%s", <-hash1, <-hash2)
}

// MultiHash

func MultiHash(in, out chan interface{}) {
	m := &sync.Mutex{}
	wg := &sync.WaitGroup{}

	for i := range in {
		wg.Add(1)
		go multiHashWorker(wg, m, out, fmt.Sprintf("%v", i))
	}
	wg.Wait()
}

func multiHashWorker(wg *sync.WaitGroup, m *sync.Mutex, out chan interface{}, data string) {
	defer wg.Done()
	wg2 := sync.WaitGroup{}
	parts := [6]string{}
	for th := 0; th <= 5; th++ {
		wg2.Add(1)
		go func(th int) {
			defer wg2.Done()
			hash := DataSignerCrc32(fmt.Sprintf("%d%s", th, data)) // crc32(th+data)
			m.Lock()
			parts[th] = hash
			m.Unlock()
		}(th)
	}
	wg2.Wait()
	out <- strings.Join(parts[:], "")
}

// CombineResults

func CombineResults(in, out chan interface{}) {
	results := make([]string, 0)

	for data := range in {
		results = append(results, fmt.Sprintf("%v", data))
	}

	sort.Strings(results)
	result := strings.Join(results, "_")
	//fmt.Println("result: ",result)
	out <- result
}
