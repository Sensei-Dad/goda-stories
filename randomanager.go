package main

import (
	"crypto/rand"
	"math/big"
)

// Random int from 0 to (X - 1)
func RandomInt(upperBound int) int {
	x, _ := rand.Int(rand.Reader, big.NewInt(int64(upperBound)))
	return int(x.Int64())
}

// Random int from 1 to X
// RollDie(X) is like rolling an X-sided die
func RollDie(numSides int) int {
	return RandomInt(numSides) + 1
}

func Clamp(num, minNum, maxNum int) int {
	if num < minNum {
		return minNum
	} else if num > maxNum {
		return maxNum
	} else {
		return num
	}
}
