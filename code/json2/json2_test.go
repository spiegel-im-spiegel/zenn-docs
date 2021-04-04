package jsonbench

import (
	"bytes"
	"encoding/json"
	"testing"

	another "github.com/goccy/go-json"
	emoji "github.com/spiegel-im-spiegel/emojis/json"
	"github.com/spiegel-im-spiegel/fetch"
)

func getMustEmojiSequenceJSON() []byte {
	u, err := fetch.URL("https://raw.githubusercontent.com/spiegel-im-spiegel/emojis/main/json/emoji-sequences.json")
	if err != nil {
		panic(err)
	}
	resp, err := fetch.New().Get(u)
	if err != nil {
		panic(err)
	}
	b, err := resp.DumpBodyAndClose()
	if err != nil {
		panic(err)
	}
	return b
}

var jsonText = getMustEmojiSequenceJSON()

func BenchmarkDecodeOrgPkg(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = json.NewDecoder(bytes.NewReader(jsonText)).Decode(&([]emoji.EmojiSequence{}))
	}
}

func BenchmarkDecodeAnotherPkg(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = another.NewDecoder(bytes.NewReader(jsonText)).Decode(&([]emoji.EmojiSequence{}))
	}
}
