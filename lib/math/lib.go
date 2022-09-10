package math

// Abs returns the absolute value of an integer
func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// Mod returns the Euclidean division remainder between 2 integers
func Mod(a, b int) int {
	m := a % b
	if m < 0 {
		m += Abs(b)
	}
	return m
}
