/*
Exercise 4.4: Write a version of rotate that operates in a single pass.
*/

package rotate

type TempInt struct {
	num   int
	valid bool
}

func Rotate(s []int, r int) {
	var size int = len(s)
	var swaps int = 0
	var temp TempInt
	if r%size == 0 {
		return
	}
	for i := 0; swaps < size; swaps++ {
		ri := (i + r) % size
		if temp.valid {
			next_temp := s[ri]
			s[ri] = temp.num
			temp.num = next_temp
		} else {
			temp.num = s[ri]
			temp.valid = true
			s[ri] = s[i]
		}
		i = ri
	}
}
