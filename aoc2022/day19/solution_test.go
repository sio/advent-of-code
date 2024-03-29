package main

import (
	"testing"

	"fmt"
	"strings"
)

var workers = map[string](func(string) string){
	"part1": part1,
	"part2": part2,
}

func TestSample(t *testing.T) {
	tests := []struct {
		worker string
		input  string
		result string
	}{
		{"part1", "sample.txt", "33"}, // input.txt: 1258 is too low
		{"part2", "sample.txt", fmt.Sprint(56 * 62)},
		{"part1", "sample2.txt", "72"},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("%s/%s", test.input, test.worker), func(t *testing.T) {
			got := strings.TrimSpace(workers[test.worker](test.input))
			want := strings.TrimSpace(test.result)
			if got != want {
				if strings.Contains(got, "\n") {
					t.Errorf("want:\n%s\n\ngot:\n%s", want, got)
				} else {
					t.Errorf("want %q, got %q", want, got)
				}
			}
		})
	}
}

func BenchmarkPart1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		part1("sample.txt")
	}
}

func BenchmarkPart2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		part2("sample.txt")
	}
}
