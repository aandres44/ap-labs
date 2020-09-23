package main

import (
	"golang.org/x/tour/wc"
	"strings"
)

func WordCount(s string) map[string]int {
	words := strings.Split(s," ")
	countMap := make(map[string]int)
	for i:= 0; i < len(words); i++ {
		countMap[words[i]] += 1
	}
	return countMap;
}

func main() {
	wc.Test(WordCount)
}