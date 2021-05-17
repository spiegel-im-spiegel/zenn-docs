---
title: "Go で簡単 Excel → CSV 変換" # 記事のタイトル
emoji: "💻" # アイキャッチとして使われる絵文字（1文字だけ）
type: "tech" # "tech" : 技術記事 / "idea" : アイデア記事
topics: ["go", "programming"] # タグ。["markdown", "rust", "aws"] のように指定する
published: true # 公開設定（true で公開）
---

いつもの小ネタで。

Excel を使えば XLS(X) ファイルを CSV テキストに変換できるけど，その度に Excel を起動するのは面倒だし，なにより UTF-8 でエクスポートすると忌々しい BOM (Byte Order Mark) が付いてくるのが嫌だったので Go で簡単なコマンドライン・ツールを組んでみる。

今回は [github.com/360EntSecGroup-Skylar/excelize][excelize] パッケージを使ってみた。こんな感じ。

```go:main.go
package main

import (
    "encoding/csv"
    "fmt"
    "io"
    "os"

    "github.com/360EntSecGroup-Skylar/excelize/v2"
)

func ExcelToCsv(w io.Writer, path string, sheetIndex int) error {
    excel, err := excelize.OpenFile(path)
    if err != nil {
        return err
    }
    rows, err := excel.Rows(excel.GetSheetName(sheetIndex))
    if err != nil {
        return err
    }
    csvw := csv.NewWriter(w)
    defer csvw.Flush()
    for rows.Next() {
        cols, err := rows.Columns()
        if err != nil {
            return err
        }
        if err := csvw.Write(cols); err != nil {
            return err
        }
    }
    return nil
}

func main() {
    if err := ExcelToCsv(os.Stdout, "./foo.xlsx", 0); err != nil {
        fmt.Fprintln(os.Stderr, err)
    }
}
```

これで foo.xlsx ファイルの最初のシートの内容を CSV 形式 UTF-8 文字エンコーディングで標準出力に出力する。大きいファイルだと CSV 出力を開始するまでにタイムラグが発生するのだが，ファイル全体をヒープに展開しでるのかなぁ。行ごとの列数が不揃いでも問題ないようだ。

[excelize] パッケージはセル単位でかなり細かい制御ができるみたいだ。また貼り付けた画像データの抽出もできるっぽい。色々試してみるといいだろう。

## 参考

https://qiita.com/hiro_nico/items/0f180f2dfc62cf2559c7
https://text.baldanders.info/release/2021/05/xls2csv/

[Go]: https://golang.org/ "The Go Programming Language"
[excelize]: https://github.com/360EntSecGroup-Skylar/excelize "360EntSecGroup-Skylar/excelize: Golang library for reading and writing Microsoft Excel™ (XLSX) files."
<!-- eof -->
