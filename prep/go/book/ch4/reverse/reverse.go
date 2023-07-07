/*
Exercise 4.3: Rewrite reverse to use an array pointer instead of a slice.
*/

package reverse

// reverse reverses a slice of ints in place.
func Reverse(s []int) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}
