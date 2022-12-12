package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
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

func ReadChars(filename string) (chars chan rune) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	chars = make(chan rune)
	go func() {
		defer file.Close()
		reader := bufio.NewReader(file)
		for {
			r, _, err := reader.ReadRune()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			chars <- r
		}
		close(chars)
	}()
	return chars
}

func execute(part func(string) string, input string, number int) {
	var result string
	result = part(input)
	var delimiter string
	if strings.Contains(result, "\n") {
		delimiter = "\n"
	}
	fmt.Printf("Part %d result: %s%s\n", number, delimiter, result)
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
		execute(part1, input, 1)
	case 2:
		execute(part2, input, 2)
	default:
		execute(part1, input, 1)
		execute(part2, input, 2)
	}
}
