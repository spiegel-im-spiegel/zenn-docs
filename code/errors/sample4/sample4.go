package main

import (
	"errors"
	"fmt"
	"os"
	"syscall"

	"github.com/hashicorp/go-multierror"
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
	paths := []string{"not-exist1.txt", "not-exist2.txt"}
	var result *multierror.Error
	for _, path := range paths {
		if err := checkFileOpen(path); err != nil {
			result = multierror.Append(result, err)
		}
	}
	if err := result.ErrorOrNil(); err != nil {
		//fmt.Fprintln(os.Stderr, err)
		var perr *os.PathError
		if errors.As(err, &perr) && errors.Is(perr, syscall.ENOENT) {
			fmt.Fprintf(os.Stderr, "\"%v\" ファイルが存在しない\n", perr.Path)
		} else {
			fmt.Fprintln(os.Stderr, "その他のエラー")
		}
	}
	// Output:
	// 2 errors occurred:
	//     * error! : open not-exist1.txt: no such file or directory
	//     * error! : open not-exist2.txt: no such file or directory
}
