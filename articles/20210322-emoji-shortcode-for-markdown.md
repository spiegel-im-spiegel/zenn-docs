---
title: "Markdown 用の絵文字コードの一覧を作ってみる" # 記事のタイトル
emoji: "💮" # アイキャッチとして使われる絵文字（1文字だけ）
type: "idea" # "tech" : 技術記事 / "idea" : アイデア記事
topics: ["go", "markdown", "emoji"] # タグ。["markdown", "rust", "aws"] のように指定する
published: true # 公開設定（ true で公開）
---

今回は軽く小ネタで。

[Go] 製 SSG の [Hugo] が v0.82.0 にアップデートしたので[リリースノート](https://github.com/gohugoio/hugo/releases/tag/v0.82.0 "Release v0.82.0 · gohugoio/hugo")を見ていたのだが， markdown 用の絵文字コードのパースって [kyokomi/emoji][emoji] パッケージを使ってるんだねぇ。

Markdown 用の絵文字コードってのは， GitHub や Slack なんかで `:smile:` と入力したら 😄 に変換されるやつ。

ん？ もしかして [kyokomi/emoji][emoji] パッケージがあれば絵文字コードの一覧が作れるんじゃね？ と思いついたのでコードを書いてみた。こんな感じ。

```go
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
```

で実際に動かしてみたら三千行以上の巨大テーブルになってしまった（笑） ここに貼り付けるわけにもいかないので Gist に貼っている。

- [Emoji Shortcode List · GitHub](https://gist.github.com/spiegel-im-spiegel/66aac732f27ad69cc8b6bd33478ecfa4)

ご笑覧あれ。

Zenn の markdown でも絵文字コードに対応してくれんもんかねぇ。

## 参考

https://text.baldanders.info/remark/2020/10/emoji-variation-and-markdown/

[Go]: https://golang.org/ "The Go Programming Language"
[Hugo]: https://gohugo.io/ "The world’s fastest framework for building websites | Hugo"
[emoji]: https://github.com/kyokomi/emoji "kyokomi/emoji: emoji terminal output for golang"
<!-- eof -->
