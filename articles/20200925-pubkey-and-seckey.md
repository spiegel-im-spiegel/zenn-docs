---
title: "公開鍵暗号の秘密鍵は絶対に渡してはならない（フリじゃないよ）" # 記事のタイトル
emoji: "🔐" # アイキャッチとして使われる絵文字（1文字だけ）
type: "idea" # "tech" : 技術記事 / "idea" : アイデア記事
topics: ["cryptography"] # タグ。["markdown", "rust", "aws"] のように指定する
published: true # 公開設定（true で公開）
---

Twitter で[こういう tweet](https://twitter.com/issei_y/status/1308749237253844992) を見かけて

@[tweet](https://twitter.com/issei_y/status/1308749237253844992)

まぁ，内容自体は明らかにネタなので笑ってスルーしようかと思ったが，続くスレッドがどうにも「？？？」なので，「ネタにマジレスwww」と嗤われることを承知で書いておく。

## 暗号とは

まずは，そもそもの話から。世の中には「暗号」も「暗合」も「隠蔽」も「符牒」も区別できない人がいるみたいなので，この記事における「暗号」を定義する。

暗号を数式で概念的に表すとこんな感じ。最初の式が「暗号化」，次の式が「復号」を表す。

$$
\begin{aligned}
  S' &= F(S,K_1)  \\
  S &= F^{-1}(S',K_2)
\end{aligned}
$$

暗号化では元のデータ $S$ に対し関数 $F$ をパラメータ $K_1$ と共にあてはめ，新たなデータ $S'$ を生成する。ちなみに，元のデータ $S$ を「平文」，新たなデータ $S'$ を「暗号文」と呼ぶ。暗号文 $S'$ から平文 $S$ が推測（解読）できないのが特徴。復号では暗号文 $S'$ に対し関数 $F^{-1}$ をパラメータ $K_2$ と共にあてはめ，平文 $S$ を復元する。

関数 $F$ および $F^{-1}$ のセットを「暗号方式」または単に「アルゴリズム」と呼び，パラメータ $K_1$ および $K_2$ を「鍵」と呼ぶ。このように，データ・アルゴリズム・鍵の3つの要素を組み合わせてデータを秘匿する技術が「暗号」である。 

よく出来た暗号は，強度（暗号の破れにくさ）が鍵の強度（「鍵長」と呼ぶ）のみに依存するよう設計されている。コンピュータシステムにおいてはアルゴリズムはプログラムコードとして記述されるため秘密にできない。一方，鍵はプログラムとは独立のデータなので，**鍵さえ知られなければ**平文を安全に管理することができる。

定義と用語はこんなところかな。

## 共通鍵暗号と公開鍵暗号

上の式の鍵 $K_1$, $K_2$ が同じ値となる暗号方式を「共通鍵暗号」と呼ぶ。「対称暗号」「秘密鍵暗号」と呼ばれることもある。

共通鍵暗号の特徴は以下の通り。

- 暗号化・復号の処理速度が速い
- 平文と暗号文とでデータサイズがほぼ同じ
- 暗号文をやり取りするもの同士で鍵を共有する必要があり，かつ第三者に知られてはならない（鍵配送問題）

一方，鍵 $K_1$, $K_2$ が異なる値となる暗号方式を「公開鍵暗号」と呼ぶ。また，暗号化 $F$ に使う鍵 $K_1$ を「公開鍵」，復号 $F^{-1}$ に使う鍵 $K_2$ を「秘密鍵」と呼ぶ。公開鍵から秘密鍵を推測できないのが特徴である（逆は可能）。

他に公開鍵暗号の特徴は以下の通り。

- 公開鍵は第三者に知られても構わない。秘密鍵は**誰にも**知られてはならない
- 暗号化・復号の処理速度が共通鍵暗号に比べて遅い
- 平文に対して暗号文のサイズが巨大になる（概ね倍以上）
- 同じ暗号強度の共通鍵暗号と比べて鍵長が巨大になる（アルゴリズムにもよるが，2倍から数十倍）

## ハイブリッド暗号

上述したように，共通鍵暗号は使い勝手のいい優れた暗号方式だが，鍵配送問題という致命的な欠点を持つ。逆に公開鍵暗号は効率も使い勝手も悪いが**秘密鍵を共有する必要がない**という一点に於いて非常に優秀である。

そこで，実際のデータ暗号の運用としては共通鍵暗号と公開鍵暗号を組み合わせた「ハイブリッド暗号」が使われる。

たとえば，以下の図は OpenPGP でメッセージ（平文）を暗号化する際の手順を示したものだ。

[![OpenPGP による暗号化](https://baldanders.info/spiegel/openpgp/hybrid-enc.svg =500x)](https://baldanders.info/spiegel/openpgp/ "わかる！ OpenPGP 暗号 — 旧コンテンツ置き場 | Baldanders.info")

復号手順は以下の図の通り。

[![OpenPGP による復号](https://baldanders.info/spiegel/openpgp/hybrid-dec.svg =500x)](https://baldanders.info/spiegel/openpgp/ "わかる！ OpenPGP 暗号 — 旧コンテンツ置き場 | Baldanders.info")

このように送信者の鍵は一切使う必要がないことが分かるだろう。

「これだと暗号文の送信者は復号できないぢゃん」と思う人もいるかもしれないが，セッション鍵を送信者と受診者の2つの公開鍵で暗号化すれば問題ない。

たとえば OpenPGP 実装のひとつである GnuPG では

```
$ echo hello, world | gpg -ea -r alice -r bob
-----BEGIN PGP MESSAGE-----

hF4DUA0A1lShMzISAQdA2offsY8f1eSp5d7jVc7u9RsXQsFjPSBXcNME4BfcgVIw
P17ibgwnl9QTNpOSUUi3877AKuy4Oblp4QkiPbNQ4mHttq/Eq2pWyWmC2fMh14QW
hF4DGCdC1cDKKBUSAQdA7SPVL/VxyoNvgWxT7Fx9oswgFSg1oJ+q/aVZjyARzF8w
D2fqYjsnfa1CKAo9uwHkIIGPLgOc3VXlH9mMr/jrxQyLDYUOlCXVIqshBhRQsx7l
0kgBrmGKFcFR4MXiNnK8Y6bJxtm3koru1FrPQUPYx8/1tbWGzyqL7b5oHXtsL8tb
wY/NB5Nl6o7oJ51Yo12mflHKx6NOM6r9ruI=
=pFvw
-----END PGP MESSAGE-----
```

のように複数の公開鍵を使ってメッセージ `hello, world` を暗号化できる。ちなみに拙作の [gpgpdump] を使って可視化すると

```
$ echo hello, world | gpg -ea -r alice -r bob | gpgpdump
Public-Key Encrypted Session Key Packet (tag 1) (94 bytes)
	Version: 3 (current)
	Key ID: 0x500d00d654a13332
	Public-key Algorithm: ECDH public key algorithm (pub 18)
	ECDH EC point (Native point format of the curve follows) (263 bits)
	symmetric key (encoded) (48 bytes)
Public-Key Encrypted Session Key Packet (tag 1) (94 bytes)
	Version: 3 (current)
	Key ID: 0x182742d5c0ca2815
	Public-key Algorithm: ECDH public key algorithm (pub 18)
	ECDH EC point (Native point format of the curve follows) (263 bits)
	symmetric key (encoded) (48 bytes)
Sym. Encrypted Integrity Protected Data Packet (tag 18) (72 bytes)
	Encrypted data (plain text + MDC SHA1(20 bytes); sym alg is specified in pub-key encrypted session key)
```

のように暗号化されたセッション鍵パケットが複数あることが分かるだろう。

ちなみに `gpg.conf` ファイルに

```
default-key alice
default-recipient-self
```

などと書いておくと，常に `alice` の鍵で暗号化してくれる。

### 【おまけ】もしも当局に「秘密鍵を渡せ！」と言われたら

GnuPG にはセッション鍵を取り出すオプションがある。

まずは通常の復号に `--show-session-key` を付加すると

```
$ cat hello.asc | gpg --show-session-key -d
gpg: 256-ビットECDH鍵, ID 182742D5C0CA2815, 日付2020-09-25に暗号化されました
      "Bob <bob@example.com>"
gpg: 256-ビットECDH鍵, ID 500D00D654A13332, 日付2020-09-25に暗号化されました
      "Alice <alice@example.com>"
gpg: session key: '9:AF823E9A36B9E4E49A2715DAD055DEE23E4169C0BFE4DAA8A7EC330582F34515'
hello, world
```

てな感じでセッション鍵 `'9:AF823...'` が取り出せる。復号時にこのセッション鍵を `--override-session-key` オプションで指定すれば

```
$ cat hello.asc | gpg --override-session-key "9:AF823E9A36B9E4E49A2715DAD055DEE23E4169C0BFE4DAA8A7EC330582F34515" -d
gpg: 256-ビットECDH鍵, ID 182742D5C0CA2815, 日付2020-09-25に暗号化されました
      "Bob <bob@example.com>"
gpg: 256-ビットECDH鍵, ID 500D00D654A13332, 日付2020-09-25に暗号化されました
      "Alice <alice@example.com>"
hello, world
```

のように，パスフレーズを訊かれることなく，秘密鍵なしで復号できる。

この機能は捜査当局等から秘密鍵を要求された際に，捜査に必要な暗号文だけを復号できるよう，セッション鍵のみを渡すためのオプションである。

## だからね

公開鍵暗号の秘密鍵は絶対に渡しちゃダメ。フリじゃないよ！

## 参考文献

続きは結城浩さんの

https://www.amazon.co.jp/dp/B015643CPE

を読むことをお勧めする。特に公開鍵が誰に所属するかの証明（certification）は公開鍵暗号および公開鍵暗号基盤の根幹に関わるので，この本を読んで理解を進めて欲しい。

あと，経路の暗号で頻出する「鍵交換」についても件の本が役に立つだろう。

### その他の関連リンク

- [GnuPG チートシート（簡易版）](./20200920-gnupg-cheat-sheet)
- [OpenPGP の電子署名は「ユーザーの身元を保証し」ない — OpenPGP の実装 | text.Baldanders.info](https://text.baldanders.info/openpgp/web-of-trust/)
- [経路の暗号化とデータの暗号化では要件が異なる — しっぽのさきっちょ | text.Baldanders.info](https://text.baldanders.info/remark/2020/07/requirement-for-encryption/)

[gpgpdump]: https://text.baldanders.info/release/gpgpdump/ "OpenPGP パケットを可視化する gpgpdump — リリース情報 | text.Baldanders.info"
