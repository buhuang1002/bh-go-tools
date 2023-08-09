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
	UnwrapReader() io.Reader
}

func UnwrapReader(r io.Reader) io.Reader {
	for {
		if wr, ok := r.(WrapReader); ok {
			r = wr.UnwrapReader()
			continue
		}

		return r
	}
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
	if len(p) > sr.rate {
		p = p[:sr.rate]
	}

	n, err := sr.r.Read(p)
	time.Sleep(time.Duration(float64(time.Second)*(float64(n)/(float64(sr.rate)))) - time.Since(t0))
	return n, err
}

func (sr *SpeedReader) UnwrapReader() io.Reader {
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

func (ir *SkipReader) UnwrapReader() io.Reader {
	return ir.r
}

var _ WrapReader = &SkipReader{}

type WrapWriter interface {
	io.Writer
	UnwrapWriter() io.Writer
}

func UnwrapWriter(w io.Writer) io.Writer {
	for {
		if ww, ok := w.(WrapWriter); ok {
			w = ww.UnwrapWriter()
			continue
		}

		return w
	}
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

func (sw *SkipWriter) UnwrapWriter() io.Writer {
	return sw.w

}

var _ WrapWriter = &SkipWriter{}

func NewLimitWriter(w io.Writer, limit int64) *LimitWriter {
	return &LimitWriter{
		w: w,
		n: limit,
	}
}

type LimitWriter struct {
	w io.Writer
	n int64
}

func (lw *LimitWriter) Write(p []byte) (int, error) {
	if lw.n == 0 {
		return 0, io.ErrShortWrite
	}

	rawLen := len(p)
	if lw.n < int64(rawLen) {
		p = p[:lw.n]
	}

	n, err := lw.w.Write(p)
	lw.n -= int64(n)

	if n == rawLen {
		return n, err
	}

	if err != nil {
		return n, err
	}

	return n, io.ErrShortWrite
}

func (lw *LimitWriter) UnwrapWriter() io.Writer {
	return lw.w
}

var _ WrapWriter = &LimitWriter{}
