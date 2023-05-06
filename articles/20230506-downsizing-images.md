---
title: "画像ファイルのサイズを縮小したい" # 記事のタイトル
emoji: "💻" # アイキャッチとして使われる絵文字（1文字だけ）
type: "tech" # "tech" : 技術記事 / "idea" : アイデア記事
topics: ["go", "programming"] # タグ。["markdown", "rust", "aws"] のように指定する
published: true # 公開設定（true で公開）
---

2023年 GW をいかがお過ごしでしょうか。私はうっかり Bluesky に手を付けてしまい，以下の公式 [Go] パッケージで遊んでみたのですが

https://github.com/bluesky-social/indigo

画像ファイルのアップロードでハマってしまいました。どうやら現状の公式 PDS では 1MB 以上の画像ファイルはサーバ側にアップロードできないみたいで（しかもアップロード失敗を返しこない）どうにか画像データを 1MB 以下に抑えようと試行錯誤してしまいました。今回はその辺の話を記しておきましょう。

さて，ここからは口調を改めて本題。

画像データのファイルサイズを縮小する手段としては以下が挙げられるだろう（不可逆圧縮になるのは諦める）。

1. JPEG 形式に変換する（特に PNG 形式に対しては効果大）
2. 画像の品質を落とす（JPEG の場合）
3. 画像のサイズを小さくする

というわけで，この方針に従って実際に画像データの縮小関数を書いてみよう。

## 完成形

とりあえず完成形はこんな感じになった。

```go:github.com/goark/toolbox/images/images.go
const (
    imageMaxSize     = 1000
    imageFileMaxSize = 1024 * 1024
)

func AjustImage(src []byte) (io.Reader, error) {
    // check file size
    if len(src) < imageFileMaxSize {
        return bytes.NewReader(src), nil
    }

    // decode image
    imgSrc, t, err := image.Decode(bytes.NewReader(src))
    if err != nil {
        return nil, errs.Wrap(err)
    }
    // convert JPEG
    quality := 90
    if t != "jpeg" {
        b, err := convertJPEG(imgSrc, quality)
        if err != nil {
            return nil, errs.Wrap(err)
        }
        if len(b) < imageFileMaxSize {
            return bytes.NewReader(b), nil
        }
    }
    // quality down
    for _, q := range []int{85, 55, 25} {
        b, err := convertJPEG(imgSrc, q)
        if err != nil {
            return nil, errs.Wrap(err)
        }
        quality = q
        if len(b) < imageFileMaxSize {
            return bytes.NewReader(b), nil
        }
    }

    // rectange of image
    rctSrc := imgSrc.Bounds()
    rate := 1.0
    if rctSrc.Dx() > rctSrc.Dy() {
        if rctSrc.Dx() > imageMaxSize {
            rate = imageMaxSize / float64(rctSrc.Dx())
        }
    } else {
        if rctSrc.Dy() > imageMaxSize {
            rate = imageMaxSize / float64(rctSrc.Dy())
        }
    }
    if rate >= 1.0 {
        return nil, errs.Wrap(ecode.ErrTooLargeImage)
    }

    // scale down
    dstX := int(float64(rctSrc.Dx()) * rate)
    dstY := int(float64(rctSrc.Dy()) * rate)
    imgDst := image.NewRGBA(image.Rect(0, 0, dstX, dstY))
    draw.BiLinear.Scale(imgDst, imgDst.Bounds(), imgSrc, rctSrc, draw.Over, nil)
    b, err := convertJPEG(imgDst, quality)
    if err != nil {
        return nil, errs.Wrap(err)
    }
    if len(b) > imageFileMaxSize {
        return nil, errs.Wrap(ecode.ErrTooLargeImage)
    }
    return bytes.NewReader(b), nil
}

func convertJPEG(src image.Image, quality int) ([]byte, error) {
    dst := &bytes.Buffer{}
    if err := jpeg.Encode(dst, src, &jpeg.Options{Quality: quality}); err != nil {
        return nil, errs.Wrap(err)
    }
    return dst.Bytes(), nil
}
```

前提として画像データは []byte 型のスライスに格納済みとする。これを AjustImage() 関数に渡すわけだ。

## image.Image 型にデコードする

詳しく見ていこう。まずはサイズを確認。 1MB 以下のデータサイズならそのまま返す。

```go
if len(src) < imageFileMaxSize {
    return bytes.NewReader(src), nil
}
```

次に画像データを [image].Image 型にデコードする。

```go
imgSrc, t, err := image.Decode(bytes.NewReader(src))
if err != nil {
    return nil, errs.Wrap(err)
}
```

## JPEG 形式に変換してみる

作成した `imgSrc` を JPEG 形式のバイナリに変換するのは簡単。以下の関数でできる。

```go
func convertJPEG(src image.Image, quality int) ([]byte, error) {
    dst := &bytes.Buffer{}
    if err := jpeg.Encode(dst, src, &jpeg.Options{Quality: quality}); err != nil {
        return nil, errs.Wrap(err)
    }
    return dst.Bytes(), nil
}
```

これを使って JPEG 以外の画像データを変換してサイズをチェックする。

```go
quality := 90
if t != "jpeg" {
    b, err := convertJPEG(imgSrc, quality)
    if err != nil {
        return nil, errs.Wrap(err)
    }
    if len(b) < imageFileMaxSize {
        return bytes.NewReader(b), nil
    }
}
```

この変換により 1MB 以下のサイズになれば，そのデータを返却している。品質は90で設定している（JPEG の品質は90以上では殆ど変わらないそうな）。

## 品質を落としてみる

続けて， JPEG データも含め，品質を 85, 55, 25 と落としていってデータサイズをチェックしていく。

```go
for _, q := range []int{85, 55, 25} {
    b, err := convertJPEG(imgSrc, q)
    if err != nil {
        return nil, errs.Wrap(err)
    }
    quality = q
    if len(b) < imageFileMaxSize {
        return bytes.NewReader(b), nil
    }
}
```

品質を 85 まで落とすとファイルサイズが劇的に減るらしいので，ここで 1MB 以下になることを期待している。あとはオマケみたいなもの（笑）

## 画像のサイズを縮小してみる

品質 25 まで落としても 1MB を超える場合は，最後の手段として画像のサイズを縮小する。

まず縦横の大きい方の値を1,000ピクセルに縮小するよう比率を計算する。

```go
rctSrc := imgSrc.Bounds()
rate := 1.0
if rctSrc.Dx() > rctSrc.Dy() {
    if rctSrc.Dx() > imageMaxSize {
        rate = imageMaxSize / float64(rctSrc.Dx())
    }
} else {
    if rctSrc.Dy() > imageMaxSize {
        rate = imageMaxSize / float64(rctSrc.Dy())
    }
}
if rate >= 1.0 {
    return nil, errs.Wrap(ecode.ErrTooLargeImage)
}
```

もとより1,000ピクセルより小さいサイズであれば諦めてエラーを返している（そこから更に縮小するのはねぇ...）

次に計算した比率でスケールダウンする。

```go
dstX := int(float64(rctSrc.Dx()) * rate)
dstY := int(float64(rctSrc.Dy()) * rate)
imgDst := image.NewRGBA(image.Rect(0, 0, dstX, dstY))
draw.BiLinear.Scale(imgDst, imgDst.Bounds(), imgSrc, rctSrc, draw.Over, nil)
b, err := convertJPEG(imgDst, quality)
if err != nil {
    return nil, errs.Wrap(err)
}
if len(b) > imageFileMaxSize {
    return nil, errs.Wrap(ecode.ErrTooLargeImage)
}
return bytes.NewReader(b), nil
```

画像サイズを縮小しても 1MB を超える場合は諦めてエラーを返している。

## 実行結果

今回書いた関数を使って大きいサイズの画像ファイルも Bluesky へどうにかアップロードできるようになった。

![](/images/downsizing-images/bird.png)

メッセージ中の URL からリンクカードを表示する際のアテンション画像もこの1MB制限に引っかかることが多いみたいなので，同様に対処して，上手く表示できるようになった。

![](/images/downsizing-images/link-card.png)

こんな泥臭い方法じゃなくてもっとスマートにできるよー，というアイデアがありましたら，ぜひ教えて下さい 🙇

なお，今回作ったツールは完全に自分用（主にバッチ処理で使う予定）に作ったものなのであしからず。 ATP (Authenticated Transfer Protocol) の機能をほぼフル実装している CLI ツールとしては mattn さんのがおすすめ。

https://github.com/mattn/bsky

ここに書かれているコードはかなり参考にさせてもらっている。ありがたや。

## 参考

https://gihyo.jp/article/2023/04/bluesky-atprotocol
https://qiita.com/miyanaga/items/a616261de490cc342d08
https://text.baldanders.info/golang/resize-image/
https://text.baldanders.info/golang/resize-image-2/



[Go]: https://go.dev/ "The Go Programming Language"
[image]: https://pkg.go.dev/image "image package - image - Go Packages"
<!-- eof -->
