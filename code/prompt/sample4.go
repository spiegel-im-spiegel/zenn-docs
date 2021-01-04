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

type Buffer struct {
	head, tail int
	buffer []string
}

func NewBuffer(size int) *Buffer {
	return &Buffer{head: 0, tail: 0, buffer: make([]string, size, size)}
}

func (b *Buffer) Len() int {
	if b == nil {
		return 0
	}
	if b.head > b.tail {
		return (b.tail + len(b.buffer)) - b.head
	}
	return b.tail - b.head
}

func (b *Buffer) At(n int) string {
	if b == nil || n >= b.Len() {
		return ""
	}
	i := (b.head + n)%len(b.buffer)
	return b.buffer[i]
}

func (b *Buffer) Add(s string) {
	if b == nil {
		return
	}
	b.buffer[b.tail] = s
	b.tail = (b.tail+1)%len(b.buffer)
	if b.head == b.tail {
		b.head = (b.head+1)%len(b.buffer)
	}
}

type History struct {
	*Buffer
}

const (
	max     = 50
	logfile = "history.log"
)

func New() (*History, error) {
	history := &History{&Buffer{head: 0, tail: 0, buffer: make([]string, max, max)}}
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

func (h *History) Save() error {
	if h == nil {
		return nil
	}
	file, err := os.Create(logfile)
	if err != nil {
		return err
	}
	defer file.Close()
	for i := 0; i < h.Len(); i++ {
		fmt.Fprintln(file, h.At(i))
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
