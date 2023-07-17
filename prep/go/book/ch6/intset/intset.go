package intset

import (
	"bytes"
	"fmt"
)

const uintSize = 32 << (^uint(0) >> 63)

// An IntSet is a set of small non-negative integers.
// Its zero value represents the empty set.
type IntSet struct {
	words []uint
}

// For benchmarking
type IntBoolSet map[int]bool

func (s1 IntBoolSet) Intersect(s2 IntBoolSet) {
	for k := range s1 {
		if _, exists := s2[k]; !exists {
			delete(s1, k)
		}
	}
}

func (s1 IntBoolSet) Union(s2 IntBoolSet) {
	for k := range s2 {
		if _, exists := s1[k]; !exists {
			s1[k] = true
		}
	}
}

// Len reports the number of elements in the set
func (s *IntSet) Len() int {
	return len(s.Elems())
}

// Remove deletes x from the set
func (s *IntSet) Remove(x int) {
	// Ignore if missing
	if !s.Has(x) {
		return
	}
	word, bit := x/uintSize, uint(x%uintSize)
	s.words[word] &^= (1 << bit)
}

// Clear removes all elements from the set
func (s *IntSet) Clear() {
	s.words = []uint{}
}

// Return a copy of the set
func (s *IntSet) Copy() *IntSet {
	var copy IntSet
	copy.AddAll(s.Elems()...)
	return &copy
}

// Create a slice containing elements of the set
func (s *IntSet) Elems() []int {
	var ret []int
	for i, word := range s.words {
		if word == 0 {
			continue
		}
		for j := 0; j < uintSize; j++ {
			if word&(1<<uint(j)) != 0 {
				ret = append(ret, uintSize*i+j)
			}
		}
	}
	return ret
}

// Has reports whether the set contains the non-negative value x.
func (s *IntSet) Has(x int) bool {
	word, bit := x/uintSize, uint(x%uintSize)
	return word < len(s.words) && s.words[word]&(1<<bit) != 0
}

// Add adds the non-negative value x to the set.
func (s *IntSet) Add(x int) {
	word, bit := x/uintSize, uint(x%uintSize)
	for word >= len(s.words) {
		s.words = append(s.words, 0)
	}
	s.words[word] |= 1 << bit
}

func (s *IntSet) AddAll(nums ...int) {
	for _, n := range nums {
		s.Add(n)
	}
}

// UnionWith sets s to the union of s and t.
func (s *IntSet) UnionWith(t *IntSet) {
	for i, tword := range t.words {
		if i < len(s.words) {
			s.words[i] |= tword
		} else {
			s.words = append(s.words, tword)
		}
	}
}

func (s *IntSet) IntersectWith(t *IntSet) {
	for i := range s.words {
		if i < len(t.words) {
			s.words[i] &= t.words[i]
		} else {
			s.words[i] = 0
		}
	}
}

func (s *IntSet) DifferenceWith(t *IntSet) {
	for i := range s.words {
		if i < len(t.words) {
			s.words[i] &^= t.words[i]
		}
	}
}

func (s *IntSet) SymmetricDifference(t *IntSet) {
	for i := range t.words {
		if i < len(s.words) {
			s.words[i] ^= t.words[i]
		} else {
			s.words = append(s.words, t.words[i])
		}
	}
}

// String returns the set as a string of the form "{1 2 3}".
func (s *IntSet) String() string {
	var buf bytes.Buffer
	buf.WriteByte('{')
	for _, n := range s.Elems() {
		if buf.Len() > len("{") {
			buf.WriteByte(' ')
		}
		fmt.Fprintf(&buf, "%d", n)
	}
	buf.WriteByte('}')
	return buf.String()
}
