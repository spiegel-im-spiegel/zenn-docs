// +build run

package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/kyokomi/emoji/v2"
)

type EmojiCode struct {
	Code, Emoji string
	Aliases     []string
}

func NewEmoji(e string, cs []string) (EmojiCode, bool) {
	if len(cs) > 0 {
		return EmojiCode{Code: emoji.NormalizeShortCode(cs[0]), Emoji: e, Aliases: cs}, true
	}
	return EmojiCode{}, false
}

func EmojiListAll() []EmojiCode {
	emojiList := []EmojiCode{}
	for e, clist := range emoji.RevCodeMap() {
		if ec, ok := NewEmoji(e, clist); ok {
			emojiList = append(emojiList, ec)
		}
	}
	sort.Slice(emojiList, func(i, j int) bool {
		return strings.Compare(emojiList[i].Code, emojiList[j].Code) < 0
	})
	return emojiList
}

func main() {
	fmt.Println("| Short Code | Graph | Aliases |")
	fmt.Println("| ---------- | :---: | ------- |")
	for _, ec := range EmojiListAll() {
		var bldr strings.Builder
		for _, c := range ec.Aliases {
			if ec.Code != c {
				bldr.WriteString(fmt.Sprintf(" `%s`", c))
			}
		}
		fmt.Printf("| `%s` | %s |%s |\n", ec.Code, ec.Emoji, bldr.String())
	}
}
