package bhio

import (
	"io"
	"time"
)

var NullReader nullReader

type nullReader struct{}

func (nullReader) Read(out []byte) (int, error) {
	for i := range out {
		out[i] = 0
	}
	return len(out), nil
}

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

type WrapReader interface {
	io.Reader
	RawReader() io.Reader
}

func NewSpeedReader(r io.Reader, rate int) *SpeedReader {
	if rate < 1 {
		panic("illegal argument")
	}

	return &SpeedReader{
		r,
		rate,
	}
}

type SpeedReader struct {
	r    io.Reader
	rate int // data size transferred per second
}

func (sr *SpeedReader) Read(p []byte) (int, error) {
	t0 := time.Now()
	n, err := sr.r.Read(p)
	time.Sleep(time.Duration(int(time.Second)*n/sr.rate) - time.Since(t0))
	return n, err
}

func (sr *SpeedReader) RawReader() io.Reader {
	return sr.r
}

var _ WrapReader = &SpeedReader{}

func NewSkipReader(r io.Reader, off int64) *SkipReader {
	return &SkipReader{
		off: off,
		r:   r,
	}
}

type SkipReader struct {
	off int64
	r   io.Reader
}

func (ir *SkipReader) Read(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}

	for ir.off > 0 {
		buf := make([]byte, 32*1024)
		for {
			if int(ir.off) < len(buf) {
				buf = buf[:ir.off]
			}

			n, err := ir.r.Read(buf)
			ir.off -= int64(n)
			if err != nil {
				return 0, err
			}

			if ir.off == 0 {
				break
			}
		}
	}

	return ir.r.Read(p)
}

func (ir *SkipReader) RawReader() io.Reader {
	return ir.r
}

var _ WrapReader = &SkipReader{}

type WrapWriter interface {
	io.Writer
	RawWriter() io.Writer
}

func NewSkipWriter(w io.Writer, skipN int64) *SkipWriter {
	return &SkipWriter{
		skipN: skipN,
		w:     w,
	}
}

type SkipWriter struct {
	skipN int64
	w     io.Writer
}

func (sw *SkipWriter) Write(p []byte) (int, error) {
	if sw.skipN > 0 {
		if int(sw.skipN) >= len(p) {
			sw.skipN -= int64(len(p))
			return len(p), nil
		}

		p = p[sw.skipN:]
	}

	n, err := sw.w.Write(p)
	if err != nil && n == 0 {
		return 0, err
	}

	n = n + int(sw.skipN)
	sw.skipN = 0
	return n, err
}

func (sw *SkipWriter) RawWriter() io.Writer {
	return sw.w
}

var _ WrapWriter = &SkipWriter{}

func NewLimitWriter(w io.Writer, limit int) *LimitWriter {
	return &LimitWriter{
		w: w,
		n: limit,
	}
}

type LimitWriter struct {
	w io.Writer
	n int
}

func (lw *LimitWriter) Write(p []byte) (int, error) {
	if lw.n == 0 {
		return 0, io.ErrShortWrite
	}

	rawLen := len(p)
	if lw.n < rawLen {
		p = p[:lw.n]
	}

	n, err := lw.w.Write(p)
	lw.n -= n

	if n == rawLen {
		return n, err
	}

	if err != nil {
		return n, err
	}

	return n, io.ErrShortWrite
}

func (lw *LimitWriter) RawWriter() io.Writer {
	return lw.w
}

var _ WrapWriter = &LimitWriter{}
