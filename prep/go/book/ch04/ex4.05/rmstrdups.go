/*
Exercise 4.5: Write an in-place function to eliminate adjacent duplicates in a []string slice.
*/

package rmstrdups

func Rmstrdups(strings []string) []string {
	var i int = 1
	for _, s := range strings {
		if s != strings[i-1] {
			strings[i] = s
			i++
		}
	}
	return strings[:i]
}
