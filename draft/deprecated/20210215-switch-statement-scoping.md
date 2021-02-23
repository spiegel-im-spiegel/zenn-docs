---
title: "Switch 文のスコープ" # 記事のタイトル
emoji: "🤔" # アイキャッチとして使われる絵文字（1文字だけ）
type: "idea" # "tech" : 技術記事 / "idea" : アイデア記事
topics: [] # タグ。["markdown", "rust", "aws"] のように指定する
published: false # 公開設定（true で公開）
---

たとえば [Go] で

```go
package main

import "fmt"

func main() {
	v := 1
	switch v {
	case 1:
		say := "yes"
		fmt.Println("say", say)
	case 2:
		say := "no"
		fmt.Println("say", say)
	default:
		say := "???"
		fmt.Println("say", say)
	}
}
```

というコードは普通に動く。変数 say がシャドウイングされているかどうかは switch 文を

```go
switch v {
case 1:
    say := "yes"
    fmt.Println("say", say)
case 2:
    fmt.Println("say", say)
default:
    fmt.Println("say", say)
}
```

と書き換えてコンパイルしてみれば分かる。コンパイルエラーになる筈である。

一方，たとえば Java で

```java
public class Sample1 {
	public static void main(String[] args) {
		int v = 1;
		switch (v) {
			case 1:
				String say = "yes";
				System.out.println("say "+ say); 
				break;
			case 2:
				String say = "no";
				System.out.println("say "+ say); 
				break;
			default:
				String say = "???";
				System.out.println("say "+ say); 
				break;
		}
	}
}
```





[Go]: https://golang.org/ "The Go Programming Language"
