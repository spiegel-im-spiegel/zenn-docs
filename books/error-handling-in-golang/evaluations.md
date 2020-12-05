---
title: "エラー評価のいろいろ"
---

エラーハンドリングを行うためにはまず何らかの形でエラーを評価する必要がある。 [Go] におけるエラー評価は大雑把に以下の4つに分けられるだろう。

1. nil との比較
2. インスタンスの同値性（equality）
3. インスタンスのボックス化解除（unboxing）
4. Error() メソッドの返り値（文字列）を解析する

以降からひとつづつ見ていこう。

## nil との比較（エラーの有無の判定）

おそらく [Go] で最もよく見かけるコードのパターンは

```go
if err != nil {
    ...
}
```

だろう。最近の（[Go] の支援機能を備えた）高機能エディタなら `iferr` と入力して [Tab] キーを押すと

```go
if err != nil {
    return
}
```

まで展開してくれる。重宝してます，ホンマ（笑）

### Error とボクシング

[前節](./basics)でも述べたように error は interface 型のひとつだが，そもそも interface 型の機能はボックス化（boxing）の一種と見なせる[^boxing1]。つまり

[^boxing1]: 念のために解説すると「ボックス化」とは，あるインスタンスを型と値をセットにして（大抵はヒープ上の）特定領域に格納することを言う。言い方を変えるとボックス化インスタンスは内部属性としてインスタンスの型と値を持っているわけだ。スマートポインタや依存の注入などは，このボックス化の仕組みと密接な関係がある。

```go
if err != nil {
    ...
}
```

は「err の中にボックス化されたエラー・インスタンスが入っているか」という評価である，と言えるわけだ。

interface 型と nil の関係については以下の拙文を参照のこと。

https://zenn.dev/spiegel/articles/20201010-ni-is-not-nil

## インスタンスの同値性[^equality1]

[^equality1]: IT 用語としての “equality” は日本語では何故か「等価性」と訳されることが多いが，等価性ならむしろ “equivalency” だよなぁ。というわけで，この辺の「用語」は混乱していて宗教論争に発展することも多い。私はそういうものに巻き込まれたくないので，この本では「equality ＝ 同値性」と定義する。ちなみに [Go] では `==` や `!=` は「値」の比較しかしない。ポインタ値の比較で同じ値であれば結果的に2つのインスタンスは「同一」であると見なせるが，やっていることはあくまで「値」の比較である。この辺も [Go] ならではのシンプルさと言えよう。

次によく見るのは

```go
if err != io.EOF {
    ...
}
```

のようなパターン。 [io].EOF は [io] 標準パッケージにおいて

```go:io/io.go
// EOF is the error returned by Read when no more input is available.
// Functions should return EOF only to signal a graceful end of input.
// If the EOF occurs unexpectedly in a structured data stream,
// the appropriate error is either ErrUnexpectedEOF or some other error
// giving more detail.
var EOF = errors.New("EOF")
```

などと定義されていて，ストリームの終端を示す EOF エラーとして広く使われている。なので，エラー・インスタンスが [io].EOF と同値であれば EOF エラーであると評価できるわけだ。

このように，あらかじめエラー・インスタンスを定義しておいて，それらとの比較を行うことで簡単にエラーの評価を行うことができる。

なお [Go] 1.13 からは [errors].Is() 関数が正式に用意されていて，上のコードは

```go
if !errors.Is(err, io.EOF) {
    ...
}
```

と置き換えることができる。むしろ今後は [errors].Is() 関数を使うことを強くお勧めする。

（[errors].Is() 関数については「構造化エラー」の節で再び紹介する）

## インスタンスのボックス化解除




[Go]: https://golang.org/ "The Go Programming Language"
[io]: https://golang.org/pkg/io/ "io - The Go Programming Language"
[errors]: https://golang.org/pkg/errors/ "errors - The Go Programming Language"
<!-- eof -->
