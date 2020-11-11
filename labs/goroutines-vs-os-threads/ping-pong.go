package main

import (
	"fmt"
	"time"
)

func main() {
	in := make(chan int)
	out := make(chan int)
	commsPerSecond := 0

	go connect(in, out)
	go connect(out, in)

	for i := 0; i < 1; i++ {
		in <- 0
		time.Sleep(time.Duration(1) * time.Second)
		x := <-in
		commsPerSecond += x
	}

	fmt.Println("Communications Per Second : ", commsPerSecond)
}

func connect(in chan int, out chan int) {
	for {
		out <- (1 + <-in)
	}
}