// +build run

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	emoji "github.com/spiegel-im-spiegel/emojis/json"
	"github.com/spiegel-im-spiegel/fetch"
)

func getEmojiSequenceJSON() ([]byte, error) {
	u, err := fetch.URL("https://raw.githubusercontent.com/spiegel-im-spiegel/emojis/main/json/emoji-sequences.json")
	if err != nil {
		return nil, err
	}
	resp, err := fetch.New().Get(u)
	if err != nil {
		return nil, err
	}
	return resp.DumpBodyAndClose()
}

func main() {
	b, err := getEmojiSequenceJSON()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	list := []emoji.EmojiSequence{}
	if err := json.NewDecoder(bytes.NewReader(b)).Decode(&list); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	fmt.Println("| Sequence | Shortcodes |")
	fmt.Println("| :------: | ---------- |")
	for _, ec := range list {
		var bldr strings.Builder
		for _, c := range ec.Shortcodes {
			bldr.WriteString(fmt.Sprintf(" `%s`", c))
		}
		fmt.Printf("| %v |%s |\n", ec.Sequence, bldr.String())
	}
}
