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

func NewEmoji(c, e string) *EmojiCode {
	return &EmojiCode{Code: emoji.NormalizeShortCode(c), Emoji: e, Aliases: []string{}}
}

func (ec *EmojiCode) Add(code ...string) *EmojiCode {
	ec.Aliases = append(ec.Aliases, code...)
	return ec
}

func EmojiListAll() []*EmojiCode {
	emojiList := []*EmojiCode{}
	for e, clist := range emoji.RevCodeMap() {
		if len(clist) > 0 {
			emojiList = append(emojiList, NewEmoji(clist[0], e).Add(clist...))
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
