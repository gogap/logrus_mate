package logrus_mate

import (
	"io"
	"os"
)

func init() {
	RegisterWriter("stdout", NewStdoutWriter)
	RegisterWriter("stderr", NewStderrWriter)
}

func NewStdoutWriter(options Options) (writer io.Writer, err error) {
	writer = os.Stdout
	return
}

func NewStderrWriter(options Options) (writer io.Writer, err error) {
	writer = os.Stderr
	return
}
