package main

import (
	"errors"
	"fmt"
	"os"
	"syscall"
)

func checkFileOpen(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("error! : %w", err)
	}
	defer file.Close()
	return nil
}

func main() {
	if err := checkFileOpen("not-exist.txt"); err != nil {
		var perr *os.PathError
		if errors.As(err, &perr) {
			switch {
			case errors.Is(err, syscall.ENOENT):
				fmt.Fprintf(os.Stderr, "\"%v\" ファイルが存在しない\n", perr.Path)
			default:
				fmt.Fprintln(os.Stderr, "その他の PathError")
			}
		} else {
			fmt.Fprintln(os.Stderr, "その他のエラー")
		}
		return
	}
}
