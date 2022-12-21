package main

import (
	"testing"
)

func TestDistances(t *testing.T) {
	tests := []struct {
		from, to string
		distance int
	}{
		{"AA", "DD", 1},
		{"DD", "BB", 5 - 2 - 1},
		{"BB", "JJ", 9 - 5 - 1},
		{"JJ", "HH", 17 - 9 - 1},
		{"EE", "HH", 21 - 17 - 1},
		{"EE", "CC", 24 - 21 - 1},
	}

	tunnels := &Graph{}
	err := tunnels.ParseFile("sample.txt")
	if err != nil {
		t.Fatal(err)
	}

	var a, b *Valve
	var got int
	var ok bool
	for _, test := range tests {
		a, ok = tunnels.Get(test.from)
		if !ok {
			t.Errorf("node not found: %s", test.from)
			continue
		}
		b, ok = tunnels.Get(test.to)
		if !ok {
			t.Errorf("node not found: %s", test.to)
			continue
		}
		got = tunnels.Distance(a, b)
		if got != test.distance {
			t.Errorf("incorrect distance between %v and %v: expected %d, got %d", a, b, test.distance, got)
			continue
		}
		got = tunnels.Distance(b, a)
		if got != test.distance {
			t.Errorf("mismatching distance in reverse direction between %v and %v: expected %d, got %d", b, a, test.distance, got)
		}
	}
}

// With my input the winning path yields a reward = 1724:
//	 [AI KB QK CJ KS CU YE]
