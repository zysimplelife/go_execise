package main

import (
	"fmt"
	"time"
)

func MinuteTicker() *time.Ticker {
	c := make(chan time.Time, 1)
	t := &time.Ticker{C: c}
	go func() {
		for {
			n := time.Now()
			if n.Second() == 0 {
				c <- n
			}
			time.Sleep(time.Second)
		}
	}()
	return t
}

func main() {
	for n := range MinuteTicker().C {
		fmt.Println("NOW: ", n)
	}
}
