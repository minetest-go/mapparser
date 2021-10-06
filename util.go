package mapparser

func getNodePos(x, y, z int) int {
	return x + (y * 16) + (z * 256)
}

func fromNodePos(pos int) (int, int, int) {
	x := pos % 16
	pos /= 16
	y := pos % 16
	pos /= 16
	z := pos % 16
	return x, y, z
}
