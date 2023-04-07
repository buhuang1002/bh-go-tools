package bhio

var NullReader nullReader

type nullReader struct{}

func (nullReader) Read(out []byte) (int, error) {
	for i := range out {
		out[i] = 0
	}
	return len(out), nil
}
