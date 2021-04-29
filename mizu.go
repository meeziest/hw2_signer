package main

import "fmt"

func main() {
	inputData := []int{0, 1, 1, 2, 3, 5, 8}
	in := make(chan interface{}, 100)
	out := make(chan interface{}, 100)
	for _, fibNum := range inputData {
		out <- fibNum
	}
	close(out)
	in = out
	out = make(chan interface{}, 100)
	SingleHash(in, out)
	fmt.Println("SingleHash: ", <-out)
	close(out)
	in = out
	out = make(chan interface{}, 100)
	MultiHash(in, out)
	fmt.Println("MultiHash: ", <-out)
	close(out)
	in = out
	out = make(chan interface{}, 100)
	CombineResults(in, out)
	fmt.Println("CombineResult: ", <-out)
}
