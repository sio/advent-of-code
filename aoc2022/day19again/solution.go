package main

import (
	"fmt"
	"log"
)

func part1(filename string) string {
	var iter LineIterator
	err := iter.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer iter.Close()
	for iter.Next() {
		fmt.Printf("> %s\n", iter.Value())
	}
	return ""
}

func part2(filename string) string {
	return ""
}
