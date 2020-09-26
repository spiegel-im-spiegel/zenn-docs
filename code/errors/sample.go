package main

import (
	"fmt"
	"os"
	"syscall"
)

func checkFileOpen(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return nil
}

// func main() {
// 	if err := checkFileOpen("not-exist.txt"); err != nil {
// 		var errPath *os.PathError
// 		if errors.As(err, &errPath) {
// 			switch {
// 			case errors.Is(errPath.Err, syscall.ENOENT):
// 				fmt.Fprintf(os.Stderr, "%v ファイルが存在しない\n", errPath.Path)
// 			default:
// 				fmt.Fprintln(os.Stderr, "その他の PathError")
// 			}
// 		} else {
// 			fmt.Fprintln(os.Stderr, "その他のエラー")
// 		}
// 		return
// 	}
// 	fmt.Println("正常終了")
// }

// func main() {
// 	if err := checkFileOpen("not-exist.txt"); err != nil {
// 		switch {
// 		case errors.Is(err, syscall.ENOENT):
// 			fmt.Fprintln(os.Stderr, "ファイルが存在しない")
// 		default:
// 			fmt.Fprintln(os.Stderr, "その他のエラー")
// 		}
// 		return
// 	}
// 	fmt.Println("正常終了")
// }

func main() {
	if err := checkFileOpen("not-exist.txt"); err != nil {
		switch e := err.(type) {
		case *os.PathError:
			if errno, ok := e.Err.(syscall.Errno); ok {
				switch errno {
				case syscall.ENOENT:
					fmt.Fprintf(os.Stderr, "%v ファイルが存在しない\n", e.Path)
				default:
					fmt.Fprintln(os.Stderr, "Errno =", errno)
				}
			} else {
				fmt.Fprintln(os.Stderr, "その他の PathError")
			}
		default:
			fmt.Fprintln(os.Stderr, "その他のエラー")
		}
		return
	}
	fmt.Println("正常終了")
}
