# 絵文字の取得

## github.com/kyokomi/emoji/v2 パッケージ

Hugo でも使われている。

```go
package main

import (
    "fmt"

    "github.com/kyokomi/emoji/v2"
)

func main() {
    fmt.Println("Hello World Emoji!")

    emoji.Println(":beer: Beer!!!")

    pizzaMessage := emoji.Sprint("I like a :pizza: and :sushi:!!")
    fmt.Println(pizzaMessage)
}
```

実行結果

```
$ go run sample1.go
Hello World Emoji!
🍺  Beer!!!
I like a 🍕  and 🍣 !!
```

https://gist.github.com/spiegel-im-spiegel/66aac732f27ad69cc8b6bd33478ecfa4
