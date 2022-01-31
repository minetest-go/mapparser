package mapparser

import (
	"fmt"
	"testing"
)

func TestToNodePos(t *testing.T) {
	pos := GetNodePos(5, 8, 12)

	x, y, z := FromNodePos(pos)
	fmt.Println(pos)

	if x != 5 {
		t.Errorf("x mismatch: %d", x)
	}

	if y != 8 {
		t.Errorf("y mismatch: %d", y)
	}

	if z != 12 {
		t.Errorf("z mismatch: %d", z)
	}
}
