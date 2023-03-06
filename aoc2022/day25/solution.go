package main

func part1(filename string) string {
	var iter LineIterator
	if err := iter.Open(filename); err != nil {
		panic(err)
	}

	var total, cursor SnafuNumber
	for iter.Next() {
		err := cursor.Parse(iter.Value())
		if err != nil {
			panic(err)
		}
		total += cursor
	}
	return total.String()
}

func part2(filename string) string {
	return ""
}
