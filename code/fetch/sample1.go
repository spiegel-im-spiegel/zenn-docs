package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/spiegel-im-spiegel/fetch"
)

func main() {
	githubUser := "spiegel-im-spiegel"
	u, err := fetch.URL("https://api.github.com/users/" + githubUser + "/gpg_keys")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	resp, err := fetch.New(
		fetch.WithHTTPClient(&http.Client{}),
		fetch.WithContext(context.Background()),
	).Get(
		u,
		fetch.WithRequestHeaderSet("Accept", "application/vnd.github.v3+json"),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer func() {
		// _, _ = io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()
	if _, err := io.Copy(os.Stdout, resp.Body); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
