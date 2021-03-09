---
title: "鬼畜上司と社畜部下（Go のチャネルであそぶ）" # 記事のタイトル
emoji: "💻" # アイキャッチとして使われる絵文字（1文字だけ）
type: "tech" # "tech" : 技術記事 / "idea" : アイデア記事
topics: ["go", "programming"] # タグ。["markdown", "rust", "aws"] のように指定する
published: true # 公開設定（true で公開）
---

突然だが，先日の[オンライン・イベント]の復習。

## 並行処理とデータ競合

書籍『[プログラミング言語Go]』では「データ競合（data race）」を以下のように定義している（9章）。

>二つのゴルーチンが同じ変数へ並行にアクセスして、そのアクセスの少なくとも一つが書き込みの場合に発生します。

そして，データ競合を避ける手段として以下の3つの方法を挙げている。

1. **変数への書き込みをしない**； immutable な構造は並行的に安全
2. **複数のゴルーチンからの変数へのアクセスを避ける**； 変数を単一のゴルーチンに閉じ込め，通信を介してデータを共有する
3. **多数のゴルーチンに変数へのアクセスを許すが，一度に一つのゴルーチンだけにアクセスさせる**； 相互排他（mutual exclusion）

最初のは Java などでよく出てくる値オブジェクト（value object）なんかを思い浮かべると分かりやすいかも知れない。3番目は，いわゆる不変式（invariant）に絡む部分で少々ややこしいので[^invariant] 今回は割愛する（私がもう少し勉強してからね）。

[^invariant]: 不変式を真面目に説明すると群論とか出てくるので，ここでは並行処理に絡めて大雑把な説明に留めるが（不正確なのはご容赦），インスタンス内の状態やインスタンス間の関係が壊れていないことを「不変式が真である」あるいは「不変式が維持されている」などと言ったりする。たとえば配列のソートを行っている最中は一時的に不変式が偽になっている。不変式が偽の状態で外部からその配列にアクセスしても内容が不定で保証されないわけ。だからソート処理全体にロックをかけて外部からアクセスさせないようにする必要がある。こういったことを状況に合わせていちいち説明するのは大変なので「不変式の真偽」という言葉で抽象化しているのだ。でも，まさに「言うは易し」で，実装コードで具体的に考えると結構ややこしかったりする。こういところが「並行（並列）処理は難しい」と思わせる理由のひとつなんだろうねぇ。

というわけで，今回は2番目の話。

## チャネルを使ったメッセージ・パッシング

[Go] ではゴルーチン（goroutine）間の通信（message-passing）にチャネル（channel）という仕組みを用意している。

チャネルは FIFO (First-In, First-Out) として機能する。さらに，チャネルに対する送受信はアトミック（atomic）であることが保証されている。複数のゴルーチンが並行に送信して値が消失したり，逆に複数ゴルーチンからの並行受信で同じ値が取り出せたりすることはないわけだ。

チャネルに何も入っていない状態（またはバッファなしのチャネル）から受信する場合，チャネルに何か入ってくるまで（またはチャネルが閉じられるまで）処理がブロックされる。

[![sanple0a.png](https://storage.googleapis.com/zenn-user-upload/b7s9914rp9cwkperud98buvi8uol)](https://storage.googleapis.com/zenn-user-upload/b7s9914rp9cwkperud98buvi8uol)

逆にチャネルのバッファがいっぱいの状態（またはバッファなしのチャネル）に送信する場合も，
チャネルからデータが取り出されるまでブロックされる。

[![sample0b.png](https://storage.googleapis.com/zenn-user-upload/tl6cfb5ud09tfrtz2fdibwpdhx62)](https://storage.googleapis.com/zenn-user-upload/tl6cfb5ud09tfrtz2fdibwpdhx62)

:::message
ブロックなしの送受信を構成することは可能。これについては後述する。
:::

他のゴルーチンとのやり取りに（共有メモリ・アクセスやメソッド経由の参照・更新などを使わず）チャネルのみを使うのであれば，並行的に安全（concurrency-safe）と言える。

## 寸劇：鬼畜上司と社畜部下

これを踏まえて，ちょっとした寸劇を考えてみた。

アクターは3人。上司1人とその部下が2人。上司はサボり魔でタスクを部下に丸投げしてとっとと家に帰りたい。部下2人は社畜で上司から仕事が降りてくるまで雛鳥のように口を開けて待っている。ある意味でよい組み合わせである（笑）

上司はチャネルの仕組みを使って簡単なタスクリストのクラスを作成した。

```go:sample1.go
// Queue: FIFO
type Queue struct {
    q chan int
}

// New: create a new instance
func New(size int) *Queue {
    return &Queue{make(chan int, size)}
}

// Add: enqueue
func (q *Queue) Add(s int) {
    q.q <- s
}

// Get: dequeue
func (q *Queue) Get() (int, bool) {
    n, ok := <-q.q
    return n, ok
}

//Complete: close channel
func (q *Queue) Complete() {
    close(q.q)
}
```

`Queue.Get()` メソッドの中身は

```go
func (q *Queue) Get() (int, bool) {
    return <-q.q
}
```

でもいいんじゃないかと思うかもしれないが，これだとコンパイルエラーになる。

```go
n, ok := <-q.q
```

は「[特殊形式（special form）](https://text.baldanders.info/golang/special-forms/ "特殊形式による式評価について | text.Baldanders.info")」なので，明示的に `(int, bool)`  で受ける必要がある。

さて，これを使って，その日に予定されているタスクを登録して部下に渡し，自分はとっとと帰宅する（笑） コードにするとこんな感じかな。

```go:sample1.go
func Manager(wg *sync.WaitGroup, tasklist []int) *Queue {
    plan := New(len(tasklist))
    wg.Add(1)
    go func() {
        defer wg.Done()
        defer plan.Complete()
        for _, n := range tasklist {
            plan.Add(n)
            log.Printf("Manager: set Task(%d)\n", n)
        }
        log.Println("Manager: return home")
    }()
    return plan
}
```

:::message
仮引数 `wg` の型がポインタ型 `*sync.WaitGroup` である点に注意。値を渡すのではなくインスタンスそのもの（への参照）を渡すわけだ。
:::

一方部下君たちの作業はこんな感じだろうか。

```go:sample1.go
const MaxWorkers = 2

func Workers(wg *sync.WaitGroup, q *Queue) {
    for i := 0; i < MaxWorkers; i++ {
        wg.Add(1)
        go func(i int) {
            defer wg.Done()
            for {
                if n, ok := q.Get(); ok {
                    log.Printf("Worker(%d): start Task(%d)\n", i, n)
                    time.Sleep(2 * time.Second) //working...
                    log.Printf("Worker(%d): end Task(%d)\n", i, n)
                } else {
                    break
                }
            }
            log.Printf("Worker(%d): return home\n", i)
        }(i + 1)
    }
}
```

タスクがなくなるまで黙々と仕事をこなす部下君たち。泣けるぜ！

チャネルがクローズされた後でも，中身が残っていれば，中身がなくなるまで受信可能である。中身がなくなったら `ok` に `false` がセットされて即時リターンとなる。

一方でクローズしたチャネルに送信すると panic になる。したがってチャネルのクローズは，基本的には，送信側の責務となる。ただし，ひとつのチャネルに対して送信ゴルーチンが複数ある場合は受信ゴルーチンを止めるための別の手立てが必要になるだろう。

では，実際にこれを実行してみよう。まず `main()` 関数はこんな感じかな。

```go:sample1.go
func main() {
    tasklist := []int{1, 2, 3, 4, 5}
    log.Println("Start...")
    var wg sync.WaitGroup
    plan := Manager(&wg, tasklist)
    Workers(&wg, plan)
    wg.Wait()
    log.Println("...End")
}
```

これを実行するとこんな感じになった。

```
$ go run sample1.go
2021/03/08 20:36:02 Start...
2021/03/08 20:36:02 Manager: set Task(1)
2021/03/08 20:36:02 Manager: set Task(2)
2021/03/08 20:36:02 Manager: set Task(3)
2021/03/08 20:36:02 Manager: set Task(4)
2021/03/08 20:36:02 Manager: set Task(5)
2021/03/08 20:36:02 Worker(2): start Task(1)
2021/03/08 20:36:02 Worker(1): start Task(2)
2021/03/08 20:36:02 Manager: return home
2021/03/08 20:36:04 Worker(2): end Task(1)
2021/03/08 20:36:04 Worker(2): start Task(3)
2021/03/08 20:36:04 Worker(1): end Task(2)
2021/03/08 20:36:04 Worker(1): start Task(4)
2021/03/08 20:36:06 Worker(1): end Task(4)
2021/03/08 20:36:06 Worker(1): start Task(5)
2021/03/08 20:36:06 Worker(2): end Task(3)
2021/03/08 20:36:06 Worker(2): return home
2021/03/08 20:36:08 Worker(1): end Task(5)
2021/03/08 20:36:08 Worker(1): return home
2021/03/08 20:36:08 ...End
```

ちょっと分かりにくいのでシーケンス図にしてみよう。こんな感じかな。

[![sample1.png](https://storage.googleapis.com/zenn-user-upload/9i7i53ettel26rcygt9bk51agevn)](https://storage.googleapis.com/zenn-user-upload/9i7i53ettel26rcygt9bk51agevn)

部下を一顧だにせず，とっとと帰宅する上司。マジ鬼畜（笑）

### チャネル送信のブロック


さて，ここで鬼畜上司にクレームが来たそうで「部下のキャパを超える仕事を一度に与えるな」ということになったらしい。そこで上司はタスクリスト（＝チャネル）のバッファを部下の数にすることで対応した。

```diff go:sample2.go
func Manager(wg *sync.WaitGroup, tasklist []int) *Queue {
-   plan := New(len(tasklist))
+   plan := New(MaxWorkers)
    wg.Add(1)
    go func() {
        defer wg.Done()
        defer plan.Complete()
        for _, n := range tasklist {
            plan.Add(n)
            log.Printf("Manager: set Task(%d)\n", n)
        }
        log.Println("Manager: return home")
    }()
    return plan
}
```

これで実行してみよう。いきなりシーケンス図で。

[![sample2.png](https://storage.googleapis.com/zenn-user-upload/7ahyiji9jqs5p7uifraoyvz1dq9w)](https://storage.googleapis.com/zenn-user-upload/7ahyiji9jqs5p7uifraoyvz1dq9w)

まぁ，上司が早めに帰ってしまうのには変わりないのだが。それよりも上司はタスクを全てセットし終わるまで部下を監視し続けるのが気に食わなかった。

そうだ。ずっと見てるんじゃなくて（自分の仕事をしながら）たまに確認するだけでいいじゃない！

### 待ちなしのチャネル送信

そこで `Queue.Add()` 関数を以下のように書き換えた。

```diff go:sample3.go
func (q *Queue) Add(s int) error {
-   q.q <- s
+   select {
+   case q.q <- s:
+       return nil
+   default:
+       return ErrTooBusy
+   }
}
```

select 文に `default` 句を付けると，どの `case` も待ち状態なら待ちなしで `default` に落ちる（チャネル受信でも同様）。これでバッファがいっぱいの場合には待ちなしでエラーが返るようになった。これを使って

```diff go:sample3.go
func Manager(wg *sync.WaitGroup, tasklist []int) *Queue {
    plan := New(MaxWorkers)
    wg.Add(1)
    go func() {
        defer wg.Done()
        defer plan.Complete()
-       for _, n := range tasklist {
-           plan.Add(n)
-           log.Printf("Manager: set Task(%d)\n", n)
+       offset := 0
+       for {
+           rest := false
+           for i := offset; i < len(tasklist); i++ {
+               offset = i
+               n := tasklist[i]
+               if err := plan.Add(n); err != nil {
+                   log.Printf("Manager: canot assign Task(%d): %v\n", n, err)
+                   rest = true
+                   break
+               } else {
+                   log.Printf("Manager: set Task(%d)\n", n)
+               }
+           }
+           if rest {
+               time.Sleep(time.Second)
+           } else {
+               break
+           }
        }
        log.Println("Manager: return home")
    }()
    return plan
}
```

てな感じに書き換えてみる。 `Queue.Add()` 関数が失敗したらいったんインターバルをおいて続きからやり直しているのがポイントである。

これでシーケンス図は

[![sample3.png](https://storage.googleapis.com/zenn-user-upload/4oxa579h820eaxk0fxd7m8hevg6f)](https://storage.googleapis.com/zenn-user-upload/4oxa579h820eaxk0fxd7m8hevg6f)

となった。

鬼畜上司は隙間時間を使って自分の仕事をすることでちょっとだけ評価が上がり，社畜な部下君たちは変わらず社畜でしたとさ。めでたしめでたし。

## ゴルーチンの優先順位

並行処理をシーケンス図で書くとどうしても「交互に実行している」ように見えてしまうのが難点だが，実際には3人のアクターを表すゴルーチンの間には優先順位はなく，完全に平等かつ並行に動く。どのゴルーチンがどんなタイミングで動作するか予測しづらいため，シビアなリアルタイム処理[^rt1] には向かないわけだ。

間接的にでもゴルーチン間に優先順位を作りたいなら，何か別の仕掛けを作り込む必要があるだろう（それでも GC なんかも考慮に入れるとかなり難しいと思うけど）。

[^rt1]: ここでいうリアルタイム処理とは「分割されたジョブを決められたタイミングで決められた期間内に完了すること」を指す。組み込みシステムではリアルタイム処理が遅滞なく行われるようジョブを設計するのが結構面倒だったりする。特にハードウェア・ブレイクしちゃう処理系はデバッグしづらいしマジ大変（笑）

[Go]: https://golang.org/ "The Go Programming Language"
[オンライン・イベント]: https://gpl-reading.connpass.com/event/204017/ "第10回『プログラミング言語Go』オンライン読書会 - connpass"
[プログラミング言語Go]: https://www.amazon.co.jp/dp/4621300253 "プログラミング言語Go (ADDISON-WESLEY PROFESSIONAL COMPUTING SERIES) | Alan A.A. Donovan, Brian W. Kernighan, 柴田 芳樹 |本 | 通販 | Amazon"

## 参考図書

https://www.amazon.co.jp/dp/4621300253
https://www.amazon.co.jp/dp/4873118468
