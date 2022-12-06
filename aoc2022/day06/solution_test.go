package main

import (
	"testing"
)

func runesChannel(input string) (output chan rune) {
	output = make(chan rune)
	go func() {
		for _, r := range input {
			output <- r
		}
		close(output)
	}()
	return output
}

func TestSamplesPart1(t *testing.T) {
	samples := map[string]int{
		"mjqjpqmgbljsphdztnvjfqwrcgsmlb":    7,
		"bvwbjplbgvbhsrlpgdmjqwftvncz":      5,
		"nppdvjthqldpwncqszvftbrmjlhg":      6,
		"nznrnfrfntjfmvfwmzdfjlvtqnbhcprsg": 10,
		"zcfzfwzzqfrljwzlrfnpqdbhtmscgvjw":  11,
	}
	var got int
	for input, expected := range samples {
		got = LocateStartOfPacket(runesChannel(input))
		if got != expected {
			t.Errorf("incorrect result for %q: expected %d, got %d", input, expected, got)
		}
	}
}
