package mapparser

// converts a vector to an integer for indexing the internal mapblock positions
func GetNodePos(x, y, z int) int {
	return x + (y * 16) + (z * 256)
}

// converts the index back to a vector
func FromNodePos(pos int) (int, int, int) {
	x := pos % 16
	pos /= 16
	y := pos % 16
	pos /= 16
	z := pos % 16
	return x, y, z
}
