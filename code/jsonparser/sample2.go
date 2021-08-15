// +build run

package main

import (
	"fmt"
	"os"

	"github.com/buger/jsonparser"
)

var jsondata = []byte(`{
  "person": {
    "name": {
      "first": "Leonid",
      "last": "Bugaev",
      "fullName": "Leonid Bugaev"
    },
    "github": {
      "handle": "buger",
      "followers": 109
    },
    "avatars": [
      { "url": "https://avatars1.githubusercontent.com/u/14009?v=3&s=460", "type": "thumbnail" }
    ]
  },
  "company": {
    "name": "Acme"
  }
}`)

func main() {
	if err := jsonparser.ObjectEach(jsondata, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		fmt.Printf("Offset: %d\n\tKey: '%s'\n\tValue: '%s'\n\tType: %s\n", offset, string(key), string(value), dataType)
		return nil
	}, "person", "name"); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
}
