---
title: "Interface 型の使いどころ【Go】" # 記事のタイトル
emoji: "🤔" # アイキャッチとして使われる絵文字（1文字だけ）
type: "idea" # "tech" : 技術記事 / "idea" : アイデア記事
topics: ["go", "programming"] # タグ。["markdown", "rust", "aws"] のように指定する
published: true # 公開設定（true で公開）
---

「[golangのコンストラクタでinterface型を返すようにする理由](https://qiita.com/asosori2/items/cec7be14c98ce4a59180)」とそこからリンクされている

https://selfnote.work/20201123/programming/how-to-use-interface-in-golang/

という記事を見て，なかなか面白いけど「理由」としてはイマイチな気がするので，この記事でも少し考えてみる。

## ”Accept Interfaces, Return Structs”

[Go] の設計指針として有名な言葉に ”Accept Interfaces, Return Structs” と言うのがある。つまり返り値としては具体的な型を返すけど，格納するインスタンスや関数の引数などは interface 型で受け入れる，というもの。

たとえば

1. 指定したファイルをオープンする
2. オープンしたファイルの内容を全て読み込む

という関数をそれぞれ作る場合

```go
func OpenFile(path string) (*os.File, error) {
    return os.Open(path)
}

func ReadAll(r io.Reader) ([]byte, error) {
    buf := bytes.Buffer{}
    _, err := buf.ReadFrom(r)
    return buf.Bytes(), err
}
```

などと書ける。これを実際に動かすには

```go
func main() {
    f, _ := OpenFile("sample.txt")
    defer f.Close()
    b, _ := ReadAll(f)
    ...
}
```

などとすればいいだろうか（エラーハンドリングをサボってます。ごめんペコン）

ここで `OpenFile()` 関数の返り値の型が `*os.File` という構造型（のポインタ）なのに対し，それを受ける `ReadAll()` 関数の引数の型は `io.Reader` という interface 型になっている点に注目である。

でも `OpenFile()` 関数の返り値を

```go
func OpenFile(path string) (io.ReadCloser, error) {
    return os.Open(path)
}
```

と interface 型の `io.ReadCloser` に書き換えても（他を書き換えることなく）全く問題なく動く。では「どちらが正しい」のだろう。

最初に紹介したリンク先の記事では「どのメソッドを『使わせる』か」はオブジェクトを返した側の責務と考えているようだ。一方 ”Accept Interfaces, Return Structs” の設計指針では「どのメソッドを『使う』か」はオブジェクトを使う側の責務と見なしている。

つまりこれは設計時における責務分担の話である，と捉えることができる。故に「どちらが正しい」のか一概には言えない，ということになる。

しかし，この一連の流れを見て面白いことに気が付かないだろうか。

### [Go] では使うメソッドを「使う側」が決定できる

実は `os.File` 構造体型と `io.ReadCloser` や `io.Reader` といった interface 型との間には明示的な「関係」は定義されていない。ただ「同じ型と名前のメソッドが定義されている」というだけである。それでもコード上は両者に関係があるように振る舞うことができる。

なので，たとえば

```go
type FileObject interface {
	Read(p []byte) (n int, err error)
	Close() error
}
```

みたいな interface 型を勝手に作って

```go
var f FileObject
f, _ := OpenFile("sample.txt")
```

などと返り値を受けても全く問題なく動く。この辺が C++ や Java などの「公称型」の型システムと根本的に異なるところで，この機能のおかげで「どのメソッドを『使う』か」をオブジェクトを使う側が決めることができるのである。

これって見方を変えると，オブジェクトの使い方を interface 型を挟んで渡す側と受ける側との合意で決定できる，ってことでもあるんだよね。つまり interface 型って「設計書」なんだよ。

これが ”Accept Interfaces, Return Structs” が示す意味で，実に [Go] らしい部分でもある。

## Interface 型を返したほうがいい場合

とはいえ常に ”Return Structs” とすべきかについて，私は必ずしも同意しない。これは個人的な見解だが interface 型を返したほうがいい場合を2つほど挙げてみる。

### 異なる型を返す可能性がある場合

分かりやすいのはエラー・ハンドリングだろう。 [Go] では

```go
// The error built-in interface type is the conventional interface for
// representing an error condition, with the nil value representing no error.
type error interface {
    Error() string
}
```

に適合する型は全てエラー・オブジェクトと見なされる。実行時エラーは，いつ・どこで・どんな事象が発生するかわからないが，とりあえず `error` 型に汎化してしまえば受ける側が適切に処理してくれる（筈）。

### レイヤ境界では Interface 型を返したほうがいい

抽象型を使うメリットのひとつはオブジェクト間の関係を「疎」にできることにある。これはプログラミング言語に依らずオブジェクト指向設計をする上で重要なポイントだろう。 “Don't Talk to Strangers” 原則というやつである。

特にチームでシステムを作っている場合，レイヤごとに人的リソースも進捗も個別に進行することが多い。前節で「interface 型は設計書」と書いたが，受け渡しする interface 型をレイヤ間であらかじめ決めておいて，それに合うようお互いにコードを書いていけば後々の結合がしやすくなる。

もっと言うと決められた interface 型に適合するモックを `xxx_test.go` とかで作れば，少なくとも単体テストは独立に行うことができる。このように [Go] はテスト駆動型の開発がしやすいよう考えられている。

## Interface 型は nil に注意

Interface 型はボックス化（boxing）の一種で nil の取り扱いに注意する必要がある。詳しくは以下の拙文をどうぞ。

https://zenn.dev/spiegel/articles/20201010-ni-is-not-nil

## Interface 型に適合するか確認するコンパイラ・ヒント

記述した型が特定の interface 型に適合しているか確認するためにこんな記述をすることがある。

```go
var _ InterfaceType = (*StructType)(nil)
```

この記述は実際にはコード化されないが `StructType` 型が interface 型の `InterfaceType` に適合しない場合はコンパイルエラーになる。コンパイラ・ヒントとして機能するわけやね。

このテクニックは割と使われるので覚えておくとよいだろう。

[Go]: https://golang.org/ "The Go Programming Language"
