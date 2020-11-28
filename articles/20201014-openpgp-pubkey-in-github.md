---
title: "GitHub に登録した OpenPGP 公開鍵を取り出す" # 記事のタイトル
emoji: "🔐" # アイキャッチとして使われる絵文字（1文字だけ）
type: "idea" # "tech" : 技術記事 / "idea" : アイデア記事
topics: ["cryptography", "openpgp", "github"] # タグ。["markdown", "rust", "aws"] のように指定する
published: true # 公開設定（true で公開）
---

[GitHub] ではコミット時やタグ付与時の電子署名を検証できる [OpenPGP] 公開鍵を登録することができるが，登録した [OpenPGP] 公開鍵を REST API を使って任意に取り出すことができる。

たとえば私が登録している公開鍵の情報はこんな感じに取り出せる[^h1]。

[^h1]: curl コマンドを使って [GitHub] REST API でアクセスする際には `-H "Accept: application/vnd.github.v3+json"` オプションを付けるのが正しいが，なくてもとりあえず問題はないっぽい。今回は GraphQL API については割愛する。

```
$ curl -s https://api.github.com/users/spiegel-im-spiegel/gpg_keys
[
  {
    "id": 1043333,
    "primary_key_id": null,
    "key_id": "B4DA3BAE7E20B81C",
    "raw_key": "-----BEGIN PGP PUBLIC KEY BLOCK-----\r\n\r\n...\r\n-----END PGP PUBLIC KEY BLOCK-----\r\n",
    "public_key": "...",
    "emails": [
      {
        "email": "spiegel.im.spiegel...",
        "verified": true
      }
    ],
    "subkeys": [
      {
        "id": 1043334,
        "primary_key_id": 1043333,
        "key_id": "4308C4946F760D3C",
        "raw_key": null,
        "public_key": "...",
        "emails": [

        ],
        "subkeys": [

        ],
        "can_sign": false,
        "can_encrypt_comms": true,
        "can_encrypt_storage": true,
        "can_certify": false,
        "created_at": "2020-10-13T22:07:23.000Z",
        "expires_at": null
      }
    ],
    "can_sign": true,
    "can_encrypt_comms": false,
    "can_encrypt_storage": false,
    "can_certify": true,
    "created_at": "2020-10-13T22:07:23.000Z",
    "expires_at": null
  }
]
```

一部 `...` で省略しているがあしからず。

公開鍵は複数登録できるので配列構造になっている。たとえば最初の鍵を取り出したいのであれば [jq] コマンドを使って

```
$ curl -s https://api.github.com/users/spiegel-im-spiegel/gpg_keys | jq .[0]
```

などとすればよい。鍵IDがあらかじめ分かっているのなら

```
$ curl -s https://api.github.com/users/spiegel-im-spiegel/gpg_keys | jq '.[]|select(.key_id=="B4DA3BAE7E20B81C")'
```

などとすることもできる。 [jq] めっさ便利！

取り出した情報のうち `raw_key` 項目の内容が [GitHub] に実際に登録した [OpenPGP] 公開鍵データだ。これを取り出すには [jq] の `-r` オプションを使って

```
$ curl -s https://api.github.com/users/spiegel-im-spiegel/gpg_keys | jq -r '.[]|select(.key_id=="B4DA3BAE7E20B81C")|.raw_key'
-----BEGIN PGP PUBLIC KEY BLOCK-----

...
-----END PGP PUBLIC KEY BLOCK-----
```

などとすればよい。この [OpenPGP] 公開鍵データはそのまま [GnuPG] 等にインポートできる。

```
$ curl -s https://api.github.com/users/spiegel-im-spiegel/gpg_keys | jq -r '.[]|select(.key_id=="B4DA3BAE7E20B81C")|.raw_key' | gpg --import
```

登録している公開鍵によっては `raw_key` 項目が `null` になっているものもあるようだ（登録時期が古いもの？）。この場合 [OpenPGP] 公開鍵として取り出すことは出来ないが，公開鍵パケットのみであれば `public_key` 項目から取り出すことは可能である。

ただし BASE64 で符号化してあるので base64 や openssl などのコマンドでバイナリデータに復号する必要がある。更に拙作の [gpgpdump] を使って取り出した公開鍵パケットを可視化することも可能である。

```
$ curl -s https://api.github.com/users/spiegel-im-spiegel/gpg_keys | jq -r '.[]|select(.key_id=="B4DA3BAE7E20B81C")|.public_key' | base64 -d | gpgpdump -u
Public-Key Packet (tag 6) (1198 bytes)
    Version: 4 (current)
    Public key creation time: 2013-04-28T10:29:43Z
    Public-key Algorithm: DSA (Digital Signature Algorithm) (pub 17)
    DSA p (3072 bits)
    DSA q (q is a prime divisor of p-1) (256 bits)
    DSA g (3070 bits)
    DSA y (= g^x mod p where x is secret) (3067 bits)
```

なお,公開鍵パケットのみでは [OpenPGP] 公開鍵として使うことは出来ないのであしからず[^pct1]。 `raw_key` 項目はないがどうしても [OpenPGP] 公開鍵を入手したいのであれば，たとえば

[^pct1]: [OpenPGP] 公開鍵は公開鍵パケット，ユーザ ID パケット，署名パケットなど複数のパケットで構成されている。公開鍵パケットのみでは鍵自体を証明することが出来ないので，少なくとも [GnuPG] では使用することが出来ない。

```
$ curl -s https://api.github.com/users/spiegel-im-spiegel/gpg_keys | jq -r .[].key_id
B4DA3BAE7E20B81C
```

といった感じに鍵IDのリストを取り出すことはできるので，以下のように

```
$ gpg --recv-keys B4DA3BAE7E20B81C
```

鍵サーバから [GnuPG] 等にインポートするしかないだろう。

## 【2020-10-14 追記】

コメントで教えてもらったのだが（情報ありがとうございます），ユーザ・プロファイルの URL に `.gpg` をくっつけたら [OpenPGP] 公開鍵データが取れるようだ。私の場合なら

```
$ curl -s https://github.com/spiegel-im-spiegel.gpg
-----BEGIN PGP PUBLIC KEY BLOCK-----

...
-----END PGP PUBLIC KEY BLOCK-----
```

で取れる。なのでこのまま

```
$ curl -s https://github.com/spiegel-im-spiegel.gpg | gpg --import
```

とするか，あるいは直接

```
$ gpg --fetch-key https://github.com/spiegel-im-spiegel.gpg
```

とすれば [GnuPG] にインポートできる。

ちなみに複数の公開鍵を登録している場合は全ての鍵データを連結した状態で取り出せる。逆にひとつも公開鍵を登録してないユーザは

```
$ curl -s https://github.com/nokeyuser.gpg
-----BEGIN PGP PUBLIC KEY BLOCK-----
Note: This user hasn't uploaded any GPG keys.


=twTO
-----END PGP PUBLIC KEY BLOCK-----
```

みたいな表示になる。

## 【2020-11-23 追記】

拙作の [gpgpdump] で [GitHub] に登録した [OpenPGP] 公開鍵を取り出して可視化できるようにした。

たとえば GitHub 上で以下のような署名を見つけたら

[![verified-signature.png](https://text.baldanders.info/release/2020/11/gpgpdump-v0_10_0-is-released/verified-signature.png)](https://text.baldanders.info/release/2020/11/gpgpdump-v0_10_0-is-released/ "gpgpdump v0.10.0 をリリースした | text.Baldanders.info")

以下のコマンドラインで公開鍵の中身を見ることができる。

```text
$ gpgpdump github spiegel-im-spiegel --keyid B4DA3BAE7E20B81C -u
Public-Key Packet (tag 6) (1198 bytes)
    Version: 4 (current)
    Public key creation time: 2013-04-28T10:29:43Z
    Public-key Algorithm: DSA (Digital Signature Algorithm) (pub 17)
    DSA p (3072 bits)
    DSA q (q is a prime divisor of p-1) (256 bits)
    DSA g (3070 bits)
    DSA y (= g^x mod p where x is secret) (3067 bits)
...
```

詳しくは以下の記事をどうぞ。

https://text.baldanders.info/release/2020/11/gpgpdump-v0_10_0-is-released/

## 参考

- [Git Commit で OpenPGP 署名を行う — OpenPGP の実装 | text.Baldanders.info](https://text.baldanders.info/openpgp/git-commit-with-openpgp-signature/)
- [jqの使い方まとめ](https://zenn.dev/syui/articles/command-json-jq)
- [GnuPG チートシート（簡易版）](https://zenn.dev/spiegel/articles/20200920-gnupg-cheat-sheet)
- [OpenPGP パケットを可視化する gpgpdump — リリース情報 | text.Baldanders.info][gpgpdump]

[GitHub]: https://github.com/
[OpenPGP]: https://tools.ietf.org/html/rfc4880 "RFC 4880 - OpenPGP Message Format"
[GnuPG]: https://www.gnupg.org/ "The GNU Privacy Guard"
[jq]: https://stedolan.github.io/jq/
[gpgpdump]: https://text.baldanders.info/release/gpgpdump/ "OpenPGP パケットを可視化する gpgpdump — リリース情報 | text.Baldanders.info"
<!-- eof -->
