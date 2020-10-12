---
title: "Zenn で “Hello World!” を埋め込んでみる" # 記事のタイトル
emoji: "💮" # アイキャッチとして使われる絵文字（1文字だけ）
type: "idea" # "tech" : 技術記事 / "idea" : アイデア記事
topics: ["programming", "markdown"] # タグ。["markdown", "rust", "aws"] のように指定する
published: false # 公開設定（true で公開）
---

まぁ，公式の「[ZennのMarkdown記法](https://zenn.dev/zenn/articles/markdown-guide)」を見れば書いてあることなのだが，この記事でちょっと試し書きしてみる。

## Markdown の Code Fence 記述でファイル名が指定できる

たえとば

> \```go:hello.go
> package main
> 
> import "fmt"
> 
> func main() {
> &nbsp;&nbsp;&nbsp;&nbsp;fmt.Println("Hello, World!")
> }
> \```

と書けば

```go:hello.go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
```

てな感じに `:` 以降の文字列を左肩に表示してくれるようになった。ブラボー







[Go]: https://golang.org/ "The Go Programming Language"
<!-- eof -->
