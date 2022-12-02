package main

import (
	"bufio"
	"log"
	"os"
)

func ReadLines(filename string) (lines chan string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	lines = make(chan string)

	go func() {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			lines <- scanner.Text()
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
		close(lines)
	}()
	return lines
}
