package main

import (
	"fmt"
	"github.com/spiegel-im-spiegel/errs"
	"os"
)

func checkFileOpen(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return errs.Wrap(
			err,
			errs.WithContext("path", path),
		)
	}
	defer file.Close()
	return nil
}

func main() {
	if err := checkFileOpen("not-exist.txt"); err != nil {
		fmt.Printf("%+v\n", err)
		return
	}
}
