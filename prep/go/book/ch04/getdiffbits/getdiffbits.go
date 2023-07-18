/*
Exercise 4.1: Write a function that counts the number of bits that are different in two SHA256 hashes.
*/

package getdiffbits

import (
	"crypto/sha256"
)

func popCount(x uint64) int {
	var c int
	for ; x > 0; c++ {
		x &= (x - 1)
	}
	return c
}

func GetDiffBits(a string, b string) int {
	var h1 = sha256.Sum256([]byte(a))
	var h2 = sha256.Sum256([]byte(b))
	var c int
	for i := range h1 {
		c += popCount(uint64(h1[i] ^ h2[i]))
	}
	return c
}
