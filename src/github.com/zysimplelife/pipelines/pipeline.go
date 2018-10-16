package main

import (
	"fmt"
	"log"
	"sync"
)

func gen(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}

func sq(done <-chan struct{}, in <-chan int, name string) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			log.Println(name, "go input", n)
			select {
			case out <- n * n:
			case <-done:
				return
			}
		}
	}()
	return out
}

func merge(done <-chan struct{}, cs ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	output := func(c <-chan int) {
		defer wg.Done()
		for n := range c {
			select {
			case out <- n:
			case <-done:
				return
			}
		}
	}

	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func main() {

	in := gen(2, 3, 4, 6)
	c1 := sq(in, "c1")
	c2 := sq(in, "c2")

	for n := range merge(c1, c2) {
		fmt.Println(n)
	}

}
