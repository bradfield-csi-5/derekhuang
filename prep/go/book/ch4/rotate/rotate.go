/*
Exercise 4.4: Write a version of rotate that operates in a single pass.
*/

package rotate

func Rotate(s []int, r int) {
	var size int = len(s)
	var swaps int = 0
	var temp = -1 // not a good starting val but not sure what it should be
	if r%size == 0 {
		return
	}
	for i := 0; swaps < size; swaps++ {
		ri := (i + r) % size
		if temp == -1 {
			temp = s[ri]
			s[ri] = s[i]
		} else {
			next_temp := s[ri]
			s[ri] = temp
			temp = next_temp
		}
		i = ri
	}
}
