package main

type PossibleMoves [][]SearchMove

func (pm *PossibleMoves) Iterator() PossibleMovesIterator {
	return PossibleMovesIterator{moves: *pm}
}

// Iterate over combinations of unique possible moves
type PossibleMovesIterator struct {
	moves     PossibleMoves
	iteration int
	ceiling   int
	value     []SearchMove
}

func (iter *PossibleMovesIterator) Initialize() {
	iter.ceiling = 1
	var i int
	for i = 0; i < len(iter.moves); i++ {
		iter.ceiling *= len(iter.moves[i])
		if iter.ceiling == 0 {
			panic("attempting to iterate over a sequence containing an empty slice")
		}
	}
}

func (iter *PossibleMovesIterator) Next() bool {
	if iter.ceiling == 0 {
		iter.Initialize()
	}

	var index int
	for {
		iter.iteration++
		iter.value = make([]SearchMove, len(iter.moves))
		index = iter.iteration - 1
		if index >= iter.ceiling {
			break
		}

		// enumerate all possible permutations
		var i, size int
		for i = 0; i < len(iter.moves); i++ {
			size = len(iter.moves[i])
			iter.value[i] = iter.moves[i][index%size]
			index = index / size
		}

		// filter out permutations with repeated values
		unique := make(map[string]bool)
		for i = 0; i < len(iter.value); i++ {
			unique[iter.value[i].Dest.Name] = true
		}
		if len(unique) == len(iter.value) {
			break
		}
	}
	return iter.iteration <= iter.ceiling
}

func (iter *PossibleMovesIterator) Value() []SearchMove {
	return iter.value
}
