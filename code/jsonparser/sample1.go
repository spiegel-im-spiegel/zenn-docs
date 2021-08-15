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
	v, err := jsonparser.GetString(jsondata, "person", "avatars", "[0]", "url")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	fmt.Println(v)
	// Output:
	// https://avatars1.githubusercontent.com/u/14009?v=3&s=460
}
