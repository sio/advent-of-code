package main

import (
	"testing"

	"fmt"
)

func TestRing(t *testing.T) {
	var numbers = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}
	sample := &Ring{}
	sample.Append(numbers...)
	fmt.Println(sample)
	if sample.Size != int64(len(numbers)) {
		t.Fatalf("expected %d items, got %d", len(numbers), sample.Size)
	}

	fmt.Println("first item:", sample.First)

	sample.First.Move(-9)
	fmt.Println(sample)
}

func TestDecryptSample(t *testing.T) {
	data := ReadCoordinates("sample.txt")
	data.Decrypt(811589153, 1)
	firstRound := []int64{0, -2434767459, 3246356612, -1623178306, 2434767459, 1623178306, 811589153}

	cursor := data.Zero
	for i := 0; i < len(firstRound); i++ {
		want := firstRound[i]
		got := cursor.Value
		if got != want {
			t.Errorf("unexpected results after first round of decryption at position %d:\nwant %d, got %d", i, want, got)
		}
		cursor = cursor.Next
	}
}

func TestDecryptFull(t *testing.T) {
	data := ReadCoordinates("sample.txt")
	data.Decrypt(811589153, 10)

	positions := map[int64]int64{
		1000: 811589153,
		2000: 2434767459,
		3000: -1623178306,
	}

	for pos, want := range positions {
		got := data.GetItem(pos)
		if got != want {
			t.Errorf("unexpected results after full decryption at position %d:\nwant %d, got %d", pos, want, got)
		}
	}
}
