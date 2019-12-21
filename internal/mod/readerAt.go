package mod

import (
	"bytes"
	"io"
)

type bufferedReaderAt struct {
	R      io.Reader
	buffer bytes.Buffer
}

func newBufferedReaderAt(r io.Reader) io.ReaderAt {
	return &bufferedReaderAt{R: r}
}

func (u *bufferedReaderAt) ReadAt(p []byte, off int64) (int, error) {
	if int(off)+len(p) > u.buffer.Len() {
		n, err := io.CopyN(&u.buffer, u.R, int64(int(off)+len(p)-u.buffer.Len()))
		if err != nil {
			buf := u.buffer.Bytes()[off : off+n]
			for i, b := range buf {
				p[i] = b
			}

			return int(n), err
		}
	}

	buf := u.buffer.Bytes()[off : int(off)+len(p)]
	for i, b := range buf {
		p[i] = b
	}

	return len(p), nil
}
