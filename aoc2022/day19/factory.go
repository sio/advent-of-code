package main

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type Factory struct {
	ID             int
	Blueprint      map[ResourceIndex]ResourcePack // map of {robot output -> robot cost}
	maxGeode       int
	maxGeodeRobots map[int]int
	debug          bool
}

func RobotOutput(output ResourceIndex) ResourcePack {
	var r ResourcePack
	if output >= 0 {
		r[output] = 1
	}
	return r
}

func (f *Factory) Debug(template string, args ...any) {
	if !f.debug {
		return
	}
	if !strings.HasSuffix(template, "\n") {
		template += "\n"
	}
	fmt.Printf(template, args...)
}

func (f *Factory) reset() {
	*f = Factory{}
}

func (f *Factory) Parse(line string) (err error) {
	f.reset() // invalidate all results of previous calculation

	var robotCost ResourcePack
	if f.Blueprint == nil {
		f.Blueprint = make(map[ResourceIndex]ResourcePack)
	}

	var words []string
	words = strings.Fields(removeSpecialCharacters(line))

	var i int
	var robotOutput, resourceType ResourceIndex
	robotOutput = Noop // all factories can do this for free
	var found bool
	for i = 0; i < len(words)-1; i++ {
		resourceType, found = ResourceName[words[i+1]] // check if the next word names a resource
		switch {

		case words[i] == "Blueprint":
			f.ID, err = strconv.Atoi(words[i+1])
			if err != nil {
				return fmt.Errorf("cannot parse blueprint ID: %w\n%s", err, line)
			}
			i++

		case words[i] == "Each":
			f.Blueprint[robotOutput] = robotCost // save previous robot

			if !found {
				return fmt.Errorf("unknown robot type (%s):\n%s", words[i+1], line)
			}
			robotOutput = resourceType
			robotCost = ResourcePack{}
			i++

		case found:
			robotCost[resourceType], err = strconv.Atoi(words[i])
			if err != nil {
				return fmt.Errorf("cannot parse cost: %s %s: %w\n%s", words[i], words[i+1], err, line)
			}
			i++

		case ignoreWords[words[i]]: // no-op

		default:
			return fmt.Errorf("parsing encountered unknown word (%s):\n%s", words[i], line)
		}
	}
	f.Blueprint[robotOutput] = robotCost // save last robot
	return nil
}

var ignoreWords = map[string]bool{
	"robot": true,
	"costs": true,
	"and":   true,
}

func removeSpecialCharacters(input string) string {
	var b strings.Builder
	var char rune
	for _, char = range input {
		if !unicode.IsLetter(char) && char != ' ' && !unicode.IsDigit(char) {
			continue
		}
		b.WriteRune(char)
	}
	return b.String()
}
