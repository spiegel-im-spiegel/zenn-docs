---
title: "GnuPG チートシート（簡易版）"
emoji: "🔐" # アイキャッチとして使われる絵文字（1文字だけ）
type: "idea" # "tech" : 技術記事 / "idea" : アイデア記事
topics: ["cryptography", "openpgp", "gnupg"] # タグ。["markdown", "rust", "aws"] のように指定する
published: true # 公開設定（true で公開）
---

:::message
この[記事は Qiita から](https://qiita.com/spiegel-im-spiegel/items/079d69282166281eb946 "GnuPG チートシート（簡易版） - Qiita")移行・再構成したものです。
:::

[GnuPG] の使い方に関する簡単な「虎の巻（cheat sheet）」を作ってみることにした。なお詳細版は以下の記事にまとめている。

http://text.baldanders.info/openpgp/gnupg-cheat-sheet/

この記事では Linux 環境を前提に書いているが Windows 等でもコマンドラインはほぼ同じなので（鍵束のあるフォルダ位置が違うが），問題なく使えると思う。

## 鍵の作成

### 対話モードでの作成

```
$ gpg --gen-key
```

[GnuPG] 2.2 系の最近のバージョンでは暗号アルゴリズムは RSA/3072bit，有効期限は作成日当日で固定されている。

暗号アルゴリズムや鍵長などを対話的に選択したい場合は `--full-gen-key` コマンドを使う。

```
$ gpg --full-gen-key
```

さらに `--expert` オプションを付けると ECC （楕円曲線暗号）が選択可能になる

```
$ gpg --full-gen-key --expert
```

ECC の取り扱いについては以下の記事を参照のこと。

http://text.baldanders.info/openpgp/using-ecc-with-gnupg/

### バッチ処理

`--gene-key` コマンドについては以下のような設定ファイルを作って `--batch` オプションを付けて起動することで対話モードを回避し，かつアルゴリズム等の詳細な指定をすることもできる。

```
$ cat alice-key.conf
Key-Type: RSA
Key-Length: 3072
Key-Usage: sign,cert
Subkey-Type: RSA
Subkey-Length: 3072
Subkey-Usage: encrypt
Name-Real: Alice
Name-Email: alice@example.com
Expire-Date: 0
Passphrase: passwd
%commit
%echo done

$ gpg --gen-key --batch alice-key.conf
gpg: 鍵3A8C942812641C17を究極的に信用するよう記録しました
gpg: ディレクトリ'/home/username/.gnupg/openpgp-revocs.d'が作成されました
gpg: 失効証明書を '/home/username/.gnupg/openpgp-revocs.d/5BD5F1F77480931FAE349B4F3A8C942812641C17.rev' に保管しました。
gpg: done

$ gpgtst --list-keys alice
pub   rsa3072 2020-09-20 [SC]
      5BD5F1F77480931FAE349B4F3A8C942812641C17
uid           [  究極  ] Alice <alice@example.com>
sub   rsa3072 2020-09-20 [E]
```

### コマンドラインでアルゴリズムを指定

```
Usage: gpg [options] --quick-generate-key user-id [algo [usage [expire]]]
```

たとえば [GnuPG] 既定のアルゴリズムで無期限の鍵を作りたいなら，こんな感じのコマンドラインになる（Pinentry でパスフレーズの入力を行う）。

```
$ gpg --quick-gen-key "Alice <alice@example.com>" default default 0
たくさんのランダム・バイトの生成が必要です。キーボードを打つ、マウスを動か
す、ディスクにアクセスするなどの他の操作を素数生成の間に行うことで、乱数生
成器に十分なエントロピーを供給する機会を与えることができます。
たくさんのランダム・バイトの生成が必要です。キーボードを打つ、マウスを動か
す、ディスクにアクセスするなどの他の操作を素数生成の間に行うことで、乱数生
成器に十分なエントロピーを供給する機会を与えることができます。
gpg: 鍵BD957725D374446Cを究極的に信用するよう記録しました
gpg: ディレクトリ'/home/username/.gnupg/openpgp-revocs.d'が作成されました
gpg: 失効証明書を '/home/username/.gnupg/openpgp-revocs.d/5852F43C85C3DF3842B9E840BD957725D374446C.rev' に保管しました。
公開鍵と秘密鍵を作成し、署名しました。

pub   rsa3072 2020-09-20 [SC]
      5852F43C85C3DF3842B9E840BD957725D374446C
uid                      Alice <alice@example.com>
sub   rsa3072 2020-09-20 [E]
```

`--quick-gen-key` コマンドでは（`default` 指定以外では）主鍵のみの作成となるので，`--quick-add-key` コマンドで副鍵を追加する。

```
Usage: gpg [options] --quick-add-key key-fingerprint [algo [usage [expire]]]
```

たとえば DSA/3072bit 主鍵を作成し，そののちに ElGamal/3072bit 副鍵を追加する場合は以下のようなコマンドラインになる。

```
$ gpg --quick-gen-key "Alice <alice@example.com>" dsa3072 default 0
...
gpg: 鍵1E2B977521CA4529を究極的に信用するよう記録しました
gpg: ディレクトリ'/home/username/.gnupg/openpgp-revocs.d'が作成されました
gpg: 失効証明書を '/home/username/.gnupg/openpgp-revocs.d/3C8FF48F3965803591FEB48C1E2B977521CA4529.rev' に保管しました。
公開鍵と秘密鍵を作成し、署名しました。

この鍵は暗号化には使用できないことに注意してください。暗号化を行うには、
"--edit-key"コマンドを使って副鍵を生成してください。
pub   dsa3072 2020-09-20 [SC]
      3C8FF48F3965803591FEB48C1E2B977521CA4529
uid                      Alice <alice@example.com>

$ gpg --quick-add-key 3C8FF48F3965803591FEB48C1E2B977521CA4529 elg3072 encr
たくさんのランダム・バイトの生成が必要です。キーボードを打つ、マウスを動か
す、ディスクにアクセスするなどの他の操作を素数生成の間に行うことで、乱数生
成器に十分なエントロピーを供給する機会を与えることができます。

$ gpg --list-keys alice
pub   dsa3072 2020-09-20 [SC]
      3C8FF48F3965803591FEB48C1E2B977521CA4529
uid           [  究極  ] Alice <alice@example.com>
sub   elg3072 2020-09-20 [E]
```

## 鍵の管理

### 鍵束内の公開鍵の検索

既に例示しているが、`--list-keys`コマンドで表示。短縮名は小文字の `-k`。

```
$ gpg -k alice
ub   dsa3072 2020-09-20 [SC]
      3C8FF48F3965803591FEB48C1E2B977521CA4529
uid           [  究極  ] Alice <alice@example.com>
sub   elg3072 2020-09-20 [E]
```

秘密鍵を検索する場合には `--list-secret-keys` コマンドを使う。短縮名は大文字の `-K`。

```
$ gpg -K alice
sec   dsa3072 2020-09-20 [SC]
      3C8FF48F3965803591FEB48C1E2B977521CA4529
uid           [  究極  ] Alice <alice@example.com>
ssb   elg3072 2020-09-20 [E]
```

### 副鍵の鍵指紋の表示

`--list-key` コマンドで主鍵の鍵指紋は表示されるが，副鍵の鍵指紋も表示したい場合は `--fingerprint` コマンドを2つ重ねる。

```
$ gpg --fingerprint --fingerprint alice
pub   dsa3072 2020-09-20 [SC]
      3C8F F48F 3965 8035 91FE  B48C 1E2B 9775 21CA 4529
uid           [  究極  ] Alice <alice@example.com>
sub   elg3072 2020-09-20 [E]
      7FDB DE0A 370F B089 4937  C70F BC59 C039 5F41 EC9F
```

### パスフレーズの変更

```
$ gpg --passwd alice
```

パスフレーズの入力は Pinentry で行う。

### 有効期限の変更

主鍵と副鍵で別々に有効期限を設定できる。

```
$ gpg --fingerprint --fingerprint alice
pub   dsa3072 2020-09-20 [SC]
      3C8F F48F 3965 8035 91FE  B48C 1E2B 9775 21CA 4529
uid           [  究極  ] Alice <alice@example.com>
sub   elg3072 2020-09-20 [E]
      7FDB DE0A 370F B089 4937  C70F BC59 C039 5F41 EC9F

gpgtst --quick-set-expire 3C8FF48F3965803591FEB48C1E2B977521CA4529 2y
$ gpg --quick-set-expire 3C8FF48F3965803591FEB48C1E2B977521CA4529 2y

$ gpg --list-keys alice
pub   dsa3072 2020-09-20 [SC] [有効期限: 2022-09-20]
      3C8FF48F3965803591FEB48C1E2B977521CA4529
uid           [  究極  ] Alice <alice@example.com>
sub   elg3072 2020-09-20 [E]

$ gpg --quick-set-expire 3C8FF48F3965803591FEB48C1E2B977521CA4529 2y 7FDBDE0A370FB0894937C70FBC59C0395F41EC9F

$ gpg --list-keys alice
pub   dsa3072 2020-09-20 [SC] [有効期限: 2022-09-20]
      3C8FF48F3965803591FEB48C1E2B977521CA4529
uid           [  究極  ] Alice <alice@example.com>
sub   elg3072 2020-09-20 [E] [有効期限: 2022-09-20]
```

前半は主鍵，後半は（主鍵に紐づく）副鍵の有効期限を変更している。

### 公開鍵をエクスポートする

```
$ gpg -a --export alice > alice-key.asc
```

### 公開鍵をインポートする

```
$ gpg --import alice-key.asc
```

Web ページ上に公開鍵のファイルを置いている場合は `--fetch-keys` コマンドで直接インポートすることもできる。

```
$ gpg --fetch-keys http://www.baldanders.info/spiegel/pubkeys/spiegel.asc
```

### 公開鍵を鍵サーバに送信する

```
$ gpg --keyserver keys.gnupg.net --send-keys 7E20B81C
```

鍵束フォルダにある `gpg.conf` ファイルに既定の鍵サーバを指定できる。

```
keyserver  keys.gnupg.net
```

主な鍵サーバは以下の通り。

- [keys.gnupg.net](http://keys.gnupg.net/ "Nebraska Wesleyan University - OpenPGP Keyserver")
- [pgp.mit.edu](https://pgp.mit.edu/ "MIT PGP Key Server")
- [pgp.nic.ad.jp](http://pgp.nic.ad.jp/ "PGP KEYSERVER")

### 公開鍵を鍵サーバから受信する

```
$ gpg --keyserver keys.gnupg.net --recv-keys 7E20B81C
```

あらかじめ鍵IDがわからない場合は `--search-keys` コマンドで以下のように検索できる。

```
$ gpg --keyserver keys.gnupg.net --search-keys alice@example.com
gpg: data source: http://hkps.pool.sks-keyservers.net:11371
(1)  Alice <alice@example.com>
    3072 bit DSA key 20EDB41D5093CA8A, 作成: 2020-06-27, 有効期限: 2022-06-27 (失効)
(2)  Gopal Moradiya <alice@example.com>
    2047 bit RSA key ABB4C0B8FA51D1CA, 作成: 2020-04-22
(3)  AliceCustom <alice_custom@example.com>
    4096 bit RSA key AFE488E4ABFB4DC5, 作成: 2019-08-22, 有効期限: 2021-08-21
(4)  example@email.cz <example@email.cz>
  Alice Vixie <alice.vixie@example.cz>
  alice.vixie@gmail.com <alice.vixie@gmail.com>
    4096 bit RSA key 41CD7069E779B2E7, 作成: 2019-04-18
(5)  alice@example.com
    2048 bit RSA key B30BCA6CF7ECB8B2, 作成: 2019-03-16
(6)  alice@example.com
    3072 bit RSA key 2CB16BC4FF043AE5, 作成: 2018-10-15 (失効)
(7)  alice <alice@example.com>
    1024 bit RSA key D06350E1BAB5DDCB, 作成: 2018-09-25, 有効期限: 2020-09-24
(8)  alice@example.com
    3072 bit RSA key 2CC0B516595D257C, 作成: 2018-09-20 (失効)
(9)  alice@example.com
    3072 bit RSA key 6D255D4FF9068C55, 作成: 2018-09-07 (失効)
(10)  alice@example.com
    3072 bit RSA key F5E8C629E71E6D26, 作成: 2018-09-02 (失効)
Keys 1-10 of 323 for "alice@example.com".  番号(s)、N)次、またはQ)中止を入力してください > 7
gpg: 鍵D06350E1BAB5DDCB: 公開鍵"alice <alice@example.com>"をインポートしました
...
```

### 公開鍵に署名する

alice の鍵で署名する場合（鍵IDを指定してもOK）。

```
$ gpg -u alice --quick-sign-key FDFA901CB9962C1FD0CB2DB3D06350E1BAB5DDCB
```

電子署名に使う既定の鍵を鍵束フォルダにある `gpg.conf` ファイルで指定することができる。

```
default-key alice
```

なお，電子署名を配布されては困る場合は `--quick-lsign-key` コマンドを使う。

## データの暗号化

### ハイブリッド暗号

alice の公開鍵で暗号化する場合（鍵IDを指定してもOK）。

```
$ echo Hello world | gpg -a -r alice -e
-----BEGIN PGP MESSAGE-----

hQMOA7xZwDlfQeyfEAv+KKzhfR/KZSZfl/pWTe7mcaA7GWydChrKkyqW2bx7I4EQ
TrQBi5UWno87Vt1TEpAgCYTIqX48N+7waI4X+jc5AjogIQW4oAp1QidEaOuINbAc
7SbaXMHqC18vVqhQCWTbmaPGj082jwlveWkHYngIQHhzoHm3JmttHwdFEvqDYmZG
dk2EeFWiHn7YpdmEoqo2aiPXPVHK7ESc6ly4K3Ris6ZpEccnUvEV/2Jy8LMs2ufP
0/N519TpuoGoCX7tzcXexgclORiWumisNnB5SGnzVRA/TIdjtaAeBhAENs91W276
Iuk/mnqopwinUGHYZLTUe2o67sR4zRRVjnz9Do6DyZcv6NqRsZmVJXeGlOfKnKMf
gxcBgTl6FopEN93r4CoaE5ccwUoeCAfym2kPDotYP79+mw8MtY2AvEYQuo0FDTJ1
iWfrBySHqmeI825RKAxFNo0pAJGV0No/oRFJVLmosIeLcgrUG0OrMfA0EcAL6sOZ
/a3hHYm8zlqyW/kyDjHLC/4tcdF29RASeYlNVuRDftH0JOOpaAt3UTVaZvrbS0qc
ZnISMnNZ04Wwi0yguinV4CE/MaR4cgrULh3nBBeQKBq0LLnRu3O4LgK7dmtRaNCF
4it3zGgW9XHPJkMgEDUH64V8IBS/hzMQWAcRG1LBzhfVg3HU5Wgh6B2rFvhU46cQ
1lxt0aoFBA9pXm1CLoMHmtHB0yu1GhO6UVBaIQvwJDxCBaZXpgVxVTPKBReelUxr
15g7i1t2bkWc+4m8MbQkVPQoSDwAST7k/fs2uOfegOMZuw6cschbt4a5cKHH+8YO
kgWkVPgegfqg/uHPhRfu1sGdY6adYEPXAJ+PLlE1+M1dqZa9uGPk/si0BFw+t7dp
1LRy89UMFkwzOW0FWBtm5qhApsdYPYbR9X/itNjMItGduw4/JYH9Yo+p5y72z1fs
ao8kfGxAlsqNfAIVXBfBgd7Esq/6o5HX4QamJU05RuDo8T3qJbn7+T+A828Lr3x/
vc8Y1znBtrwWniLV/5U+b3TSRwFshXP6Frv3s8YRmbMNFs2RO2bSlzVoqOZhruT5
EOlcpfHNbZ83pza6DQO54QFfUrbJIxITg+hB8+gE2wrDCnuxjhEl2gnO
=MVRu
-----END PGP MESSAGE-----
```

`--recipient` オプション（短縮名は `-r`）は複数指定できる。

鍵束フォルダにある `gpg.conf` ファイルで常に使用する公開鍵を指定することもできる（自分の鍵でも復号できるようにしておくのが目的）。

```
default-key alice
default-recipient-self
```

### セッション鍵のみで暗号化する

公開鍵を使わずにセッション鍵のみで暗号化することもできる。

```
$ echo Hello world | gpg -a -c
-----BEGIN PGP MESSAGE-----

jA0ECQMC8X4WEaHlqLj/0kEBpy0D0gBXMwFBFofR/QuW8WRqidUg06YlJOxh/YNK
5ErL1jwpDcQmOy00IMs32Zhf33Q3hvowiiLroTPVVGF1lw==
=D4vU
-----END PGP MESSAGE-----
```

この場合，コマンド起動時に Pinentry でパスフレーズの入力を要求され，パスフレーズからセッション鍵を生成して暗号化を行う。したがって，何らかの方法で暗号データの受け手とパスフレーズを共有する必要がある。

## 暗号データの復号

復号には `--decrypt` コマンド（短縮名は `-d`）を使う。
Pinentry によるパスフレーズの入力がある。

```
$ echo Hello world | gpg -r alice -e -a > alice-enc.asc

$ gpg -d alice-enc.asc
gpg: 3072-ビットELG鍵, ID BC59C0395F41EC9F, 日付2020-09-20に暗号化されました
      "Alice <alice@example.com>"
Hello world
```

`--output` オプション（短縮名は `-o`）で復号データの出力先ファイルを指定することもできる。
Windows の場合は安全のため（リダイレクトを使うのではなく）出力先ファイルを指定することをお勧めする。

```
$ gpg -o out.txt -d alice-enc.asc
gpg: 3072-ビットELG鍵, ID BC59C0395F41EC9F, 日付2020-09-20に暗号化されました
      "Alice <alice@example.com>"

$ cat out.txt
Hello world
```

セッション鍵のみで暗号化した場合も同じコマンドで復号できる。

```
$ echo Hello world | gpg -a -c > alice-sym-enc.asc

$ gpg -d alice-sym-enc.asc
gpg: AES256暗号化済みデータ
gpg: 1 個のパスフレーズで暗号化
Hello world
```

## データへの電子署名と検証

### クリア署名

クリア署名のコマンドは `--clear-sign`。署名者の指定には `--local-user` オプション（短縮名は `-u`）を使う（鍵ID指定でも可）。

```
$ echo Hello world | gpg -u alice --clear-sign
-----BEGIN PGP SIGNED MESSAGE-----
Hash: SHA256

Hello world
-----BEGIN PGP SIGNATURE-----

iIgEAREIADAWIQQ8j/SPOWWANZH+tIweK5d1IcpFKQUCX2c2bRIcYWxpY2VAZXhh
bXBsZS5jb20ACgkQHiuXdSHKRSmHMwD/etR+6ThaaDqmS0Z9Qmrx6q9z/WpOhH47
tIYnJpONNe8A/2taHKwG2cxlOst6gOud0cF4kdM1Bvf1Vx6YaPJXMfRO
=j0IN
-----END PGP SIGNATURE-----
```

署名検証にも `--decrypt` コマンドが使える。

```
$ echo Hello world | gpg -u alice --clear-sign | gpg -d
Hello world
gpg: 2020年09月20日 20時02分06秒 JSTに施された署名
gpg:                DSA鍵3C8FF48F3965803591FEB48C1E2B977521CA4529を使用
gpg:                発行者"alice@example.com"
gpg: "Alice <alice@example.com>"からの正しい署名 [究極]
```

### 分離署名

まず署名対象のファイルを用意する。

```
$ echo Hello world> hello.txt

$ cat hello.txt
Hello world
```

このファイルを配布した際に途中で改竄がないかどうか知りたい。
こういう場合は「分離署名」にする。
分離署名のコマンドは `--detach-sign`。短縮名は `-b`。

```
$ gpg -u alice -b hello.txt
```

この結果 `hello.txt` ファイルと同じ場所に `hello.txt.sig` ファイルが作成される。この`hello.txt` ファイルと `hello.txt.sig` ファイルをセットで配布するのである。どちらかのファイルが改竄されていれば署名の検証が NG になるはずである。

```
$ gpg --verify hello.txt.sig
gpg: 署名されたデータが'hello.txt'にあると想定します
gpg: 2020年09月20日 20時03分17秒 JSTに施された署名
gpg:                DSA鍵3C8FF48F3965803591FEB48C1E2B977521CA4529を使用
gpg:                発行者"alice@example.com"
gpg: "Alice <alice@example.com>"からの正しい署名 [究極]
```

署名対象のファイルがテキスト・ファイルの場合は `--textmode` オプションを付けて電子署名を行ったほうが安全。

```
$ gpg -u alice --textmode -b hello.txt
```

### 署名データに署名対象のデータを埋め込む

```text
$ gpg -u alice -s hello.txt
```

作成された `hello.txt.gpg` からデータの抽出と検証を行う。

```text
$ gpg -d hello.txt.gpg
Hello world
gpg: 2020年09月20日 20時05分15秒 JSTに施された署名
gpg:                DSA鍵3C8FF48F3965803591FEB48C1E2B977521CA4529を使用
gpg:                発行者"alice@example.com"
gpg: "Alice <alice@example.com>"からの正しい署名 [究極]
```

ちなみに `hello.txt.gpg` の中身を拙作の [gpgpdump] で見ると，こんな構成になっている（宣伝w）。

```
$ gpgpdump -f hello.txt.gpg -u
Compressed Data Packet (tag 8) (173 bytes)
    Compression Algorithm: ZIP <RFC1951> (comp 1)
    Compressed data (172 bytes)
    One-Pass Signature Packet (tag 4) (13 bytes)
        Version: 3 (current)
        Signiture Type: Signature of a binary document (0x00)
        Hash Algorithm: SHA2-256 (hash 8)
        Public-key Algorithm: DSA (Digital Signature Algorithm) (pub 17)
        Key ID: 0x1e2b977521ca4529
        Encrypted session key: other than one pass signature (flag 0x01)
    Literal Data Packet (tag 11) (27 bytes)
        Literal data format: b (binary)
        File name: hello.txt
        Creation time: 2020-09-20T11:05:15Z
        Literal data (12 bytes)
    Signature Packet (tag 2) (136 bytes)
        Version: 4 (current)
        Signiture Type: Signature of a binary document (0x00)
        Public-key Algorithm: DSA (Digital Signature Algorithm) (pub 17)
        Hash Algorithm: SHA2-256 (hash 8)
        Hashed Subpacket (48 bytes)
            Issuer Fingerprint (sub 33) (21 bytes)
                Version: 4 (need 20 octets length)
                Fingerprint (20 bytes)
                    3c 8f f4 8f 39 65 80 35 91 fe b4 8c 1e 2b 97 75 21 ca 45 29
            Signature Creation Time (sub 2): 2020-09-20T11:05:15Z
            Signer's User ID (sub 28): alice@example.com
        Unhashed Subpacket (10 bytes)
            Issuer (sub 16): 0x1e2b977521ca4529
        Hash left 2 bytes
            86 76
        DSA value r (254 bits)
        DSA value s (255 bits)
```

この方式の何が嬉しいかというと，対象のデータを埋め込んだ署名データをまるっと暗号化することで，暗号化と電子署名を同時に行うことができるのである。

### 暗号化と電子署名を同時に行う

たとえば Alice の鍵で署名して Bob の公開鍵で暗号化する。

```
$ gpg -u alice -r bob -se hello.txt
```

このとき生成される `hello.txt.gpg` ファイルの中身は以下の通り。

```
$ gpgpdump -f hello.txt.gpg -u
Public-Key Encrypted Session Key Packet (tag 1) (94 bytes)
    Version: 3 (current)
    Key ID: 0x35f796fd6f77fe98
    Public-key Algorithm: ECDH public key algorithm (pub 18)
    ECDH EC point (Native point format of the curve follows) (263 bits)
    symmetric key (encoded) (48 bytes)
Sym. Encrypted Integrity Protected Data Packet (tag 18) (221 bytes)
    Encrypted data (plain text + MDC SHA1(20 bytes); sym alg is specified in pub-key encrypted session key)
```

`hello.txt.gpg` ファイルから暗号データを抽出して復号と署名検証を行うには以下の通り。

```
$ gpg -d hello.txt.gpg
gpg: 256-ビットECDH鍵, ID 35F796FD6F77FE98, 日付2020-09-20に暗号化されました
      "bob <bob@example.com>"
Hello world
gpg: 2020年09月20日 20時16分31秒 JSTに施された署名
gpg:                DSA鍵3C8FF48F3965803591FEB48C1E2B977521CA4529を使用
gpg:                発行者"alice@example.com"
gpg: "Alice <alice@example.com>"からの正しい署名 [究極]
```


セッション鍵のみで暗号化する場合も同様にできる。

```
$ gpg -u alice -sc hello.txt
gpg: AES256暗号化を使用します

$ gpg -d hello.txt.gpg
gpg: AES256暗号化済みデータ
gpg: 1 個のパスフレーズで暗号化
Hello world
gpg: 2020年09月20日 20時27分18秒 JSTに施された署名
gpg:                DSA鍵3C8FF48F3965803591FEB48C1E2B977521CA4529を使用
gpg:                発行者"alice@example.com"
gpg: "Alice <alice@example.com>"からの正しい署名 [究極]
```

## 鍵の失効

鍵を作成する際に鍵束フォルダ以下の `openpgp-revocs.d` フォルダに失効証明書が作成される。中身はこんな感じ。

```
これは失効証明書でこちらのOpenPGP鍵に対するものです:

pub   dsa3072 2020-09-20 [SC]
      3C8FF48F3965803591FEB48C1E2B977521CA4529
uid          Alice <alice@example.com>

失効証明書は "殺すスイッチ" のようなもので、鍵がそれ以上使えない
ように公に宣言するものです。一度発行されると、そのような失効証明書は
撤回することはできません。

秘密鍵のコンプロマイズや紛失の場合、これを使ってこの鍵を失効させます。
しかし、秘密鍵がまだアクセス可能である場合、新しい失効証明書を生成し、
失効の理由をつける方がよいでしょう。詳細は、GnuPGマニュアルのgpgコマンド "--generate-revocation"の記述をご覧ください。

このファイルを誤って使うのを避けるため、以下ではコロンが5つのダッシュ
の前に挿入されます。この失効証明書をインポートして公開する前に、テク
スト・エディタでこのコロンを削除してください。

:-----BEGIN PGP PUBLIC KEY BLOCK-----
Comment: This is a revocation certificate

iHcEIBEIACAWIQQ8j/SPOWWANZH+tIweK5d1IcpFKQUCX2cijgIdAAAKCRAeK5d1
IcpFKSKlAP9l7iK0rhgFbSDuF93XYxihjQfU729OkaW8d7BY5wuMhgD4yN7Of9cl
qpdsGbOpRh5eEypTacAPHoxruFvFm0z06Q==
=ZQsi
-----END PGP PUBLIC KEY BLOCK-----
```

このファイルをインポートすることで鍵が失効される。
なお失効証明書を使用の際には

```
:-----BEGIN PGP PUBLIC KEY BLOCK-----
```

の先頭のコロン（`:`）を削除して使うこと。

```
$ gpg --import openpgp-revocs.d/3C8FF48F3965803591FEB48C1E2B977521CA4529.rev
gpg: 鍵1E2B977521CA4529:"Alice <alice@example.com>"失効証明書をインポートしました
gpg: 処理数の合計: 1
gpg:    新しい鍵の失効: 1

$ gpg -k alice
pub   dsa3072 2020-09-20 [SC] [失効: 2020-09-20]
      3C8FF48F3965803591FEB48C1E2B977521CA4529
uid           [  失効  ] Alice <alice@example.com>

$ gpg -a --export alice > alice-rev.asc
```

**失効した公開鍵を配布するのを忘れずに！**

失効証明書は新たに作成できる（失効理由付き）。

```
$ gpg --gen-revoke alice

sec  dsa3072/1E2B977521CA4529 2020-09-20 Alice <alice@example.com>

この鍵に対する失効証明書を作成しますか? (y/N) y
失効の理由を選択してください:
  0 = 理由は指定されていません
  1 = 鍵(の信頼性)が損なわれています
  2 = 鍵がとりかわっています
  3 = 鍵はもはや使われていません
  Q = キャンセル
(ここではたぶん1を選びたいでしょう)
あなたの決定は? 1
予備の説明を入力。空行で終了:
> 
失効理由: 鍵(の信頼性)が損なわれています
(説明はありません)
よろしいですか? (y/N) y
ASCII外装出力を強制します。
-----BEGIN PGP PUBLIC KEY BLOCK-----
Comment: This is a revocation certificate

iHgEIBEIACAWIQQ8j/SPOWWANZH+tIweK5d1IcpFKQUCX2c+MgIdAgAKCRAeK5d1
IcpFKYZwAP9d1ghgxvC7MXsUgiRcu6s91/msOFyeAfkUocgMm2P9wQD+LCGWpWmf
cirRrb34mT8g97AIkcjni0gRb0sew4Vb1gg=
=VHdm
-----END PGP PUBLIC KEY BLOCK-----
失効証明書を作成しました。

みつからないように隠せるような媒体に移してください。もし_悪者_がこの証明書への
アクセスを得ると、あなたの鍵を使えなくすることができます。
媒体が読出し不能になった場合に備えて、この証明書を印刷して保管するのが賢明です。
しかし、ご注意ください。あなたのマシンの印字システムは、他の人がアクセスできる
場所にデータをおくことがあります!
```

[OpenPGP]: http://tools.ietf.org/html/rfc4880 "RFC 4880 - OpenPGP Message Format"
[GnuPG]: https://gnupg.org/ "The GNU Privacy Guard"
[gpgpdump]: https://text.baldanders.info/release/gpgpdump/ "OpenPGP パケットを可視化する gpgpdump — リリース情報 | text.Baldanders.info"
