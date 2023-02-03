package main

import (
	"fmt"
	"log"
)

func part1(filename string) string {
	var iter CharIterator
	err := iter.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer iter.Close()
	for iter.Next() {
		fmt.Printf("> %c\n", iter.Value())
	}
	return ""
}

func part2(filename string) string {
	return ""
}
