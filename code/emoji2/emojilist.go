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
	elist := []EmojiCode{}
	for c, e := range emoji.CodeMap() {
		elist = append(elist, newEmoji(c, e))
	}
	sort.Slice(elist, func(i, j int) bool {
		return strings.Compare(elist[i].Emoji, elist[j].Emoji) < 0
	})

	emojiList := []*NormalizeEmojiCode{}
	var nec *NormalizeEmojiCode
	for _, ec := range elist {
		if !nec.Add(ec) {
			nec = newNormalizeEmoji(ec)
			emojiList = append(emojiList, nec)
		}
	}

	for i := 0; i < len(emojiList); i++ {
		if len(emojiList[i].EmojiList) > 1 {
			list := emojiList[i].EmojiList
			sort.Slice(list, func(n, m int) bool {
				return strings.Compare(list[n].Code, list[m].Code) < 0
			})
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
