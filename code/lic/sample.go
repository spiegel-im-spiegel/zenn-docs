package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/google/licenseclassifier/v2/assets"
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, os.ErrInvalid)
		return
	}
	file, err := os.Open(args[0])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer file.Close()

	c, err := assets.DefaultClassifier()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	res, err := c.MatchFrom(file)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	if len(res.Matches) == 0 {
		fmt.Fprintln(os.Stderr, args[0], "is not license file.")
		return
	}
	for _, m := range res.Matches {
		fmt.Println(m.MatchType, m.Name)
	}
}
