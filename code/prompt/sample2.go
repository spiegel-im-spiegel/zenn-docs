//go:build run
// +build run

package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/nyaosorg/go-readline-ny"
	"github.com/nyaosorg/go-readline-ny/simplehistory"
)

func Reverse(r []rune) []rune {
	if len(r) > 1 {
		for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
			r[i], r[j] = r[j], r[i]
		}
	}
	return r
}

func errPrint(err error) string {
	if err == nil {
		return ""
	}
	switch {
	case errors.Is(err, readline.CtrlC):
		return "処理を中断します"
	case errors.Is(err, io.EOF):
		return "処理を終了します"
	default:
		return err.Error()
	}
}

func main() {
	history := simplehistory.New()
	editor := readline.Editor{
		Prompt:  func() (int, error) { return fmt.Print("> ") },
		History: history,
	}
	fmt.Println("Input Ctrl+D to stop.")
	for {
		//input
		text, err := editor.ReadLine(context.Background())
		if err != nil {
			fmt.Fprintln(os.Stderr, errPrint(err))
			return
		}
		//output
		fmt.Println(string(Reverse([]rune(text))))
		//add history
		history.Add(text)
	}
	return
}
