// +build run

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func get(rawurl string) ([]byte, error) {
	resp, err := http.Get(rawurl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func main() {
	_, err := get("http://foo.bar/")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
}
