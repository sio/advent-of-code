package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
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
	cpuprofile := flag.String("cpuprofile", "", "write cpu profile to `file`")
	memprofile := flag.String("memprofile", "", "write memory profile to `file`")

	input := flag.String("input", "input.txt", "input data for today's challenge")
	part := flag.Int("part", 0, "puzzle part")
	flag.Parse()
	if flag.NArg() > 1 {
		log.Fatalf("unparsed command arguments left: %v", flag.NArg())
	}
	if flag.NArg() == 1 {
		*input = flag.Args()[0]
	}

	if *cpuprofile != "" {
		log.Printf("Writing CPU profile to %s", *cpuprofile)
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}
	if *memprofile != "" {
		log.Printf("Writing memory profile to %s", *cpuprofile)
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		runtime.GC()    // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
	}

	switch *part {
	case 1:
		execute(part1, *input, 1)
	case 2:
		execute(part2, *input, 2)
	default:
		execute(part1, *input, 1)
		execute(part2, *input, 2)
	}
}
