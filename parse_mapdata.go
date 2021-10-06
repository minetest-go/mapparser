package mapparser

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"strconv"
)

func parseMapdata(mapblock *MapBlock, data []byte) (int, error) {
	r := bytes.NewReader(data)

	cr := new(CountedReader)
	cr.Reader = r

	z, err := zlib.NewReader(cr)
	if err != nil {
		return 0, err
	}

	defer z.Close()

	buf := new(bytes.Buffer)
	io.Copy(buf, z)

	if buf.Len() != 16384 {
		return 0, errors.New("Mapdata length invalid: " + strconv.Itoa(buf.Len()))
	}

	rawdata := buf.Bytes()

	mapd := MapData{
		ContentId: make([]int, 4096),
		Param1:    make([]int, 4096),
		Param2:    make([]int, 4096),
	}
	mapblock.Mapdata = &mapd

	fmt.Println(len(rawdata))
	for i := 0; i < 4096; i++ {
		fmt.Println(i)
		mapd.ContentId[i] = int(binary.BigEndian.Uint16(rawdata[i*2:]))
		mapd.Param1[i] = int(rawdata[(4096*2)+i])
		mapd.Param2[i] = int(rawdata[(4096*3)+i])
	}

	return cr.Count, nil
}
