package main

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
)

func checkFileOpen(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return errors.WithStack(err)
		//return errors.Wrapf(err, "open error (%s)", path)
	}
	defer file.Close()
	return nil
}

func main() {
	if err := checkFileOpen("not-exist.txt"); err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		return
	}
}
