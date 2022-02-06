package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/minetest-go/mapparser"
)

// sqlite3 -header -csv map.sqlite "select pos, hex(data) from blocks > map.csv"
func main() {
	file, err := os.Open("map.csv")
	if err != nil {
		panic(err)
	}
	sc := bufio.NewScanner(file)
	line_num := 0
	for sc.Scan() {
		line := sc.Text()
		parts := strings.Split(line, ",")
		if len(parts) != 2 {
			panic(fmt.Errorf("invalid format @ %d", line_num))
		}

		if len(parts[1])%2 != 0 {
			panic(fmt.Errorf("invalid hex count @ %d, len: %d", line_num, len(parts[1])))
		}

		pos, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			panic(err)
		}

		b := make([]byte, len(parts[1])/2)
		for i := 0; i < len(parts[1]); i += 2 {
			num, err := strconv.ParseUint(parts[1][i:i+2], 16, 32)
			if err != nil {
				panic(fmt.Errorf("error @ %d: %s", line_num, err.Error()))
			}
			b[i/2] = byte(num)
		}

		mb, err := mapparser.Parse(b)
		if err != nil {
			panic(err)
		}

		x, y, z := PlainToCoord(pos)
		txt, err := json.Marshal(mb)
		if err != nil {
			panic(err)
		}

		fmt.Printf("%d;%d;%d;%s\n", x, y, z, txt)

		line_num++
	}
}
