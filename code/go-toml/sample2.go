// +build run

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/pelletier/go-toml"
	"github.com/pelletier/go-toml/query"
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, os.ErrInvalid)
		return
	}

	info, err := toml.LoadReader(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	results, err := query.CompileAndExecute(args[0], info)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	for _, item := range results.Values() {
		fmt.Println(item)
	}
}
