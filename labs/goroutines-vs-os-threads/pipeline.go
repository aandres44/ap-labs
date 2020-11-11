package main

import (
	"fmt"
	"time"
)

var maxStages int

func main() {
	in := make(chan int)
	out := make(chan int)

	go connect(in, out, 0, 1000000)
	start := time.Now()
	in <- 1

	fmt.Println("Maximum number of pipeline stages   : ", maxStages)
	fmt.Println("Time to transit trough the pipeline : ", time.Since(start))
}

func connect(in <-chan int, out chan int, index int, connections int) {
	if index != connections {
		maxStages++
		connect(in, out, index+1, connections)
		out <- <-in
	}
}
