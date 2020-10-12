package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/flynn/json5"
)

func main() {
	elms := make(map[string]interface{})
	if err := json5.NewDecoder(os.Stdin).Decode(&elms); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	if err := json.NewEncoder(os.Stdout).Encode(elms); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
