package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

type FileObject interface {
	Read(p []byte) (n int, err error)
	Close() error
}

func OpenFile(path string) (*os.File, error) {
	return os.Open(path)
}

func ReadAll(r io.Reader) ([]byte, error) {
	buf := bytes.Buffer{}
	_, err := buf.ReadFrom(r)
	return buf.Bytes(), err
}

func main() {
	var f FileObject
	f, _ = OpenFile("go.mod")
	defer f.Close()
	b, err := ReadAll(f)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	_, _ = io.Copy(os.Stdout, bytes.NewReader(b))
}
