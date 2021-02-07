package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func get2(rawurl string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, rawurl, nil)
	if err != nil {
		return nil, err
	}
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func main() {
	_, err := get2("http://foo.bar/")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
}
