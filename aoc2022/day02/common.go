package main

import (
	"bufio"
	"flag"
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

func main() {
	part := flag.Int("part", 0, "puzzle part")
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}
	input := flag.Args()[0]
	switch *part {
	case 1:
		part1(input)
	case 2:
		part2(input)
	default:
		part1(input)
		part2(input)
	}
}
