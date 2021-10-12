package mapparser

import (
	"bytes"
	"compress/zlib"
	"io"

	"github.com/klauspost/compress/zstd"
)

func decompress_zlib(data []byte) ([]byte, int, error) {
	r := bytes.NewReader(data)

	cr := new(CountedReader)
	cr.Reader = r

	z, err := zlib.NewReader(cr)
	if err != nil {
		return nil, 0, err
	}

	defer z.Close()

	buf := new(bytes.Buffer)
	io.Copy(buf, z)

	return buf.Bytes(), cr.Count, nil
}

func decompress_zstd(data []byte) ([]byte, int, error) {
	r := bytes.NewReader(data)

	cr := new(CountedReader)
	cr.Reader = r

	z, err := zstd.NewReader(cr)
	if err != nil {
		return nil, 0, err
	}

	defer z.Close()

	buf := new(bytes.Buffer)
	io.Copy(buf, z)

	return buf.Bytes(), cr.Count, nil
}
