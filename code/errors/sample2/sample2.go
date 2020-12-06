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
		switch {
		case errors.Is(err, syscall.ENOENT):
			fmt.Fprintln(os.Stderr, "ファイルが存在しない")
		default:
			fmt.Fprintln(os.Stderr, "その他のエラー")
		}
		return
	}
}
