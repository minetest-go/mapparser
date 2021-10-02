package mapparser

func getNodePos(x, y, z int) int {
	return x + (y * 16) + (z * 256)
}
