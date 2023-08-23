package bhio

import (
	"errors"
	"io"
	"sync"
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

func NewLimitWriter(w io.Writer, limit int64, endSymbol error) *LimitWriter {
	if endSymbol == nil {
		endSymbol = io.ErrShortWrite
	}

	return &LimitWriter{
		w:         w,
		n:         limit,
		endSymbol: endSymbol,
	}
}

type LimitWriter struct {
	w         io.Writer
	n         int64
	endSymbol error
}

func (lw *LimitWriter) Write(p []byte) (int, error) {
	if lw.n == 0 {
		return 0, lw.endSymbol
	}

	if lw.n < int64(len(p)) {
		p = p[:lw.n]
	}

	n, err := lw.w.Write(p)
	lw.n -= int64(n)

	if err == nil && lw.n == 0 {
		return n, lw.endSymbol
	}

	return n, err
}

func (lw *LimitWriter) UnwrapWriter() io.Writer {
	return lw.w
}

var _ WrapWriter = &LimitWriter{}

func NewBufferWriter(w io.Writer, bufLen int) *BufferWriter {
	return &BufferWriter{
		w:        w,
		readyToW: make([]byte, bufLen),
		cache:    make([]byte, bufLen),
	}
}

type BufferWriter struct {
	w        io.Writer
	wErr     error
	readyToW []byte
	cache    []byte
	wOffset  int
	w_m      sync.Mutex
}

func (bw *BufferWriter) Write(p []byte) (int, error) {
	var (
		wn int
		n  int
	)

	for {
		if bw.wErr != nil {
			return wn, bw.wErr
		}

		n = copy(bw.cache[bw.wOffset:], p[wn:])
		wn += n
		bw.wOffset += n

		if bw.wOffset == len(bw.cache) {
			bw.asyncFlush()
		}

		if wn >= len(p) {
			break
		}
	}

	return wn, nil
}

func (bw *BufferWriter) ReadFrom(r io.Reader) (int64, error) {
	var (
		wn  int
		n   int
		err error
	)

	for {
		if bw.wErr != nil {
			return int64(wn), bw.wErr
		}

		n, err = r.Read(bw.cache[bw.wOffset:])
		wn += n
		bw.wOffset += n

		if bw.wOffset == len(bw.cache) {
			bw.asyncFlush()
		}

		if err != nil {
			if errors.Is(err, io.EOF) {
				return int64(wn), nil
			}

			return int64(wn), err
		}
	}
}

func (bw *BufferWriter) DirectW(w func([]byte) error, needN int) error {
	if needN > len(bw.cache) {
		return io.ErrShortWrite
	}

	if bw.wErr != nil {
		return bw.wErr
	}

	if len(bw.cache)-bw.wOffset < needN {
		bw.asyncFlush()
	}

	err := w(bw.cache[bw.wOffset : bw.wOffset+needN])
	if err != nil {
		return err
	}

	bw.wOffset += needN
	if bw.wOffset == len(bw.cache) {
		bw.asyncFlush()
	}

	return nil
}

func (bw *BufferWriter) asyncFlush() {
	bw.flush(false)
}

func (bw *BufferWriter) flush(sync bool) {
	bw.w_m.Lock()

	if bw.wOffset == 0 {
		bw.w_m.Unlock()
		return
	}

	bw.readyToW, bw.cache = bw.cache[:bw.wOffset], bw.readyToW[:cap(bw.readyToW)]
	bw.wOffset = 0

	toWrite := func() {
		defer bw.w_m.Unlock()
		n, err := bw.w.Write(bw.readyToW)
		if err != nil {
			bw.wErr = err
			return
		}

		if n != len(bw.readyToW) {
			bw.wErr = io.ErrShortWrite
			return
		}

		return
	}

	if sync {
		toWrite()
	} else {
		go func() {
			toWrite()
		}()
	}
}

func (bw *BufferWriter) Sync() error {
	if bw.wErr != nil {
		return bw.wErr
	}

	bw.flush(true)
	return bw.wErr
}

func (bw *BufferWriter) Close() error {
	bw.cache = nil
	bw.readyToW = nil
	if closer, ok := bw.w.(io.Closer); ok {
		return closer.Close()
	}

	return nil
}

func (bw *BufferWriter) UnwrapWriter() io.Writer {
	return bw.w
}

var _ WrapWriter = &BufferWriter{}
