package bh_io

import "io"

type RepeatReader struct {
	data []byte
	i    int
}

func NewRepeatReader(data []byte) io.Reader {
	return &RepeatReader{
		data: data,
	}
}

func (rr *RepeatReader) Read(p []byte) (n int, err error) {
	for i := range p {
		p[i] = rr.data[rr.i]
		rr.i++
		if rr.i == len(rr.data) {
			rr.i = 0
		}
	}
	return len(p), nil
}
