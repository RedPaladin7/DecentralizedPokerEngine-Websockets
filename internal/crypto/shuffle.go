package crypto

import (
	"crypto/rand"
	"math/big"
)

func ShuffleDeck(deck [][]byte) [][]byte {
	n := len(deck)
	shuffled := make([][]byte, n)
	copy(shuffled, deck)

	for i := n-1; i > 0; i-- {
		jBig, err := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		if err != nil {
			continue 
		}
		j := int(jBig.Int64())
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	}
	return shuffled
}

func ApplyPermutation(deck [][]byte, permutation []int) [][]byte {
	if len(deck) != len(permutation) {
		return deck
	}

	shuffled := make([][]byte, len(deck))
	for i, idx := range permutation {
		shuffled[i] = deck[idx]
	}

	return shuffled
}

func VerifyShuffle(original, shuffled [][]byte) bool {
	if len(original) != len(shuffled) {
		return false
	}

	changedCount := 0
	for i := range original {
		if !bytesEqual(original[i], shuffled[i]) {
			changedCount++
		}
	}
	// checks if atleast 80% of the cards are in different position
	return changedCount >= len(original)*4/5
}

func bytesEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}