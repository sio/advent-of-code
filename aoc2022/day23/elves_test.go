package main

import (
	"testing"

	"fmt"
	"os"
	"strings"
)

func TestLoad(t *testing.T) {
	filename := "sample.txt"

	var got, want string

	raw, err := os.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	want = strings.TrimSpace(string(raw))

	elves := &ElfGroup{}
	elves.Load(filename)

	got = strings.TrimSpace(fmt.Sprint(elves))
	if got != want {
		for i := 0; i < len(want); i++ {
			if got[i] != want[i] {
				t.Logf("first mismatch at byte %d: want %c, got %c", i, want[i], got[i])
				break
			}
		}
		t.Errorf(
			"loaded data does not match input\nwant (%d bytes):\n%s\n\ngot (%d bytes):\n%s",
			len(want), want,
			len(got), got,
		)
	}
}

func TestPlay(t *testing.T) {
	filename := "sample.txt"

	var got, want string
	elves := &ElfGroup{}
	elves.Load(filename)
	elves.Play(10)
	got = strings.TrimSpace(fmt.Sprint(elves))

	want = strings.TrimSpace(`
......#.....
..........#.
.#.#..#.....
.....#......
..#.....#..#
#......##...
....##......
.#........#.
...#.#..#...
............
...#..#..#..`)

	if got != want {
		for i := 0; i < len(want); i++ {
			if got[i] != want[i] {
				t.Logf("first mismatch at byte %d: want %c, got %c", i, want[i], got[i])
				break
			}
		}
		t.Errorf(
			"loaded data does not match input\nwant (%d bytes):\n%s\n\ngot (%d bytes):\n%s",
			len(want), want,
			len(got), got,
		)
	}
}
