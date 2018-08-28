package main

import "time"
import "fmt"

func main() {
	//timer1 := time.NewTimer(2 * time.Second)
	//<-timer1.C
	//fmt.Println("Time 1 expired")

	ticker := time.NewTicker(5 * time.Second)
	for t := range ticker.C {
            fmt.Println("Tick arrived ", t)
        func() {
		    time.Sleep(10 * time.Second)
            fmt.Println("done with function at", t)
        }()
	}
}
