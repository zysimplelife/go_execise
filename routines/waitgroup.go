package main

import (
	"fmt"
	"sync"
)

func main() {
	var wait sync.WaitGroup
	wait.Add(5)
	for i := 0; i < 5; i++ {
		go func(id int) {
			fmt.Printf("ID:%d: Hello goroutines!\n", id)
			wait.Done()
		}(i)
	}
	wait.Wait()
}
