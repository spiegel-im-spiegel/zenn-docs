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
}

func newEmoji(c, e string) EmojiCode {
	return EmojiCode{Code: c, Emoji: e}
}

type NormalizeEmojiCode struct {
	Code      string
	EmojiList []EmojiCode
}

func newNormalizeEmoji(ec EmojiCode) *NormalizeEmojiCode {
	return &NormalizeEmojiCode{Code: emoji.NormalizeShortCode(ec.Code), EmojiList: []EmojiCode{ec}}
}

func (nec *NormalizeEmojiCode) Add(ec EmojiCode) bool {
	if nec == nil {
		return false
	}
	norm := emoji.NormalizeShortCode(ec.Code)
	if nec.Code != norm {
		return false
	}
	nec.EmojiList = append(nec.EmojiList, ec)
	return true
}

func EmojiListAll() []*NormalizeEmojiCode {
	emojiList := []*NormalizeEmojiCode{}
	for e, clist := range emoji.RevCodeMap() {
		if len(clist) > 0 {
			nec := newNormalizeEmoji(newEmoji(clist[0], e))
			for i := 1; i < len(clist); i++ {
				nec.Add(newEmoji(clist[i], e))
			}
			emojiList = append(emojiList, nec)
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
	for _, eclist := range EmojiListAll() {
		var e string
		var bldr strings.Builder
		for _, ec := range eclist.EmojiList {
			if ec.Code == eclist.Code {
				e = ec.Emoji
			} else {
				bldr.WriteString(fmt.Sprintf(" `%s`", ec.Code))
			}
		}
		fmt.Printf("| `%s` | %s |%s |\n", eclist.Code, e, bldr.String())
	}
}
