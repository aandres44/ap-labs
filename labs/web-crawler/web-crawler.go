// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 241.

// Crawl2 crawls web links starting with the command-line arguments.
//
// This version uses a buffered channel as a counting semaphore
// to limit the number of concurrent calls to links.Extract.
//
// Crawl3 adds support for depth limiting.
//
package main

import (
	"fmt"
	"log"
	"os"
	"flag"
	"gopl.io/ch5/links"
)

//!+sema
// tokens is a counting semaphore used to
// enforce a limit of 20 concurrent requests.
var tokens = make(chan struct{}, 20)

type Node struct {
	url string
	depth int
}

func crawl(node Node, fileName string, depth int) []Node {
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil{
		log.Println(err)
	}

	if _, err := file.WriteString(node.url + "\n"); err != nil{
		log.Println(err)
	}

	fmt.Println(node.url)
	file.Close()

	if node.depth < depth {
		tokens <- struct{}{} // acquire a token
		list, err := links.Extract(node.url)
		temp := make([]Node, len(list))
		for i, url := range list{
			temp[i] = Node{url: url, depth: node.depth + 1}
		}

		<-tokens // release the token

		if err != nil {
			log.Print(err)
		}

		return temp
	}
	
	return []Node{}
}

//!-sema

//!+
func main() {
	if len(os.Args) != 4 {
		fmt.Println("Error: wrong set of parameters.\n Usage: ./web-crawler -depth=<depth> -results=<results file> <url>")
		os.Exit(1)
	}

	worklist := make(chan []Node)
	var n int // number of pending sends to worklist
	initialUrl := os.Args[3:]
	depthPtr := flag.Int("depth", 1, "Depth")
	fileNamePtr := flag.String("results", "results.txt", "Result file")
	flag.Parse()

	urls := make([]Node, len(initialUrl))
	for i, url := range initialUrl {
		urls[i] = Node{url: url, depth: 0}
	}

	// Start with the command-line arguments.
	n++
	go func() { worklist <- urls }()

	// Crawl the web concurrently.
	seen := make(map[string]bool)
	for ; n > 0; n-- {
		list := <-worklist
		for _, link := range list {
			if !seen[link.url] {
				seen[link.url] = true
				n++
				go func(link Node) {
					worklist <- crawl(link, *fileNamePtr, *depthPtr)
				}(link)
			}
		}
	}
}

//!-