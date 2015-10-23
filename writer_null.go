package logrus_mate

import (
	"io"
)

func init() {
	RegisterWriter("null", NewNullWriter)
}

type NullWriter struct {
}

func (w *NullWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func NewNullWriter(options Options) (writer io.Writer, err error) {
	writer = new(NullWriter)
	return
}
