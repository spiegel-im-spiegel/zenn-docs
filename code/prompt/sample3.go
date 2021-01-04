// +build run

package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/zetamatta/go-readline-ny"
)

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

type History struct {
	buffer []string
}

var _ readline.IHistory = (*History)(nil)

const (
	max     = 50
	logfile = "history.log"
)

func New() (*History, error) {
	history := &History{buffer: []string{}}
	file, err := os.Open(logfile)
	if err != nil {
		return history, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		history.Add(scanner.Text())
	}
	return history, scanner.Err()
}

func (h *History) Len() int {
	if h == nil {
		return 0
	}
	return len(h.buffer)
}

func (h *History) At(n int) string {
	if h == nil || h.Len() <= n {
		return ""
	}
	return h.buffer[n]
}

func (h *History) Add(s string) {
	if h == nil || len(s) == 0 {
		return
	}
	if n := h.Len(); n < 1 {
		h.buffer = append(h.buffer, s)

	} else if h.buffer[n-1] != s {
		h.buffer = append(h.buffer, s)
	}
	if n := h.Len(); n > max {
		h.buffer = h.buffer[n-max:]
	}
}

func (h *History) Save() error {
	if h == nil {
		return nil
	}
	file, err := os.Create(logfile)
	if err != nil {
		return err
	}
	defer file.Close()
	for _, s := range h.buffer {
		fmt.Fprintln(file, s)
	}
	return nil
}

func Reverse(r []rune) []rune {
	if len(r) > 1 {
		for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
			r[i], r[j] = r[j], r[i]
		}
	}
	return r
}

func main() {
	history, err := New()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		//continue
	}
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
			break
		}
		//output
		fmt.Println(string(Reverse([]rune(text))))
		//add history
		history.Add(text)
	}
	if err := history.Save(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	return
}
