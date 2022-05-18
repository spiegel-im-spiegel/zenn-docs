---
title: "apt-key が非推奨になったので" # 記事のタイトル
emoji: "💮" # アイキャッチとして使われる絵文字（1文字だけ）
type: "idea" # "tech" : 技術記事 / "idea" : アイデア記事
topics: ["linux", "ubuntu", "apt"] # タグ。["markdown", "rust", "aws"] のように指定する
published: true # 公開設定（true で公開）
---

本家 [Debian] を使っている人は既にご存知だとは思うが，パッケージ管理ツールの APT でパッケージへの電子署名用 OpenPGP 公開鍵を管理する apt-key コマンドが非推奨になっている。更に [Debian] 12 では apt-key コマンドは削除されるらしい。

Debian 系ディストリビューションである [Ubuntu] も当然これに倣うのだが

```
$ sudo apt update

...

W: https://download.docker.com/linux/ubuntu/dists/jammy/InRelease: Key is stored in legacy trusted.gpg keyring (/etc/apt/trusted.gpg), see the DEPRECATION section in apt-key(8) for details.
```

みたいなワーニングは [Ubuntu] 21.10 までは出てなかった気がするので放置していた。でも，まぁ，そろそろ対応しておく必要があるだろう。幸い apt 2.4 以降[^ver] ではサードパーティの公開鍵の取り扱いについてそれほど難しくないようだ。

[^ver]: [Ubuntu] 22.04 LTS にアップグレードした時点では apt 2.4.5 になっていた。

この記事では，[Ubuntu] 22.04 LTS 以降で [Docker] Engine の APT インストール手順を例にして公開鍵のインポートやリポジトリの登録手順について紹介してみる。

## これまでのインストール手順

今は消えているみたいだが， “[Install Docker Engine on Ubuntu]” ページの古い記述では

```
$ curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
```

で公開鍵を登録し

```
$ sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
```

でリポジトリを登録する。この状態で

```text
$ sudo apt update
$ sudo apt install docker-ce docker-ce-cli containerd.io docker-compose-plugin
```

でインストールできた。

apt-key コマンドでインポートされたサードパーティの OpenPGP 公開鍵は /etc/apt/trusted.gpg 鍵束ファイルに格納されている。このファイルの中身は apt-key list コマンドで見ることができる。

```
$ apt-key list
Warning: apt-key is deprecated. Manage keyring files in trusted.gpg.d instead (see apt-key(8)).
/etc/apt/trusted.gpg
--------------------
pub   rsa4096 2017-02-22 [SCEA]
      9DC8 5822 9FC7 DD38 854A  E2D8 8D81 803C 0EBF CD88
uid           [  不明  ] Docker Release (CE deb) <docker@docker.com>
sub   rsa4096 2017-02-22 [S]
...
```

おおう。早速ワーニングが出てるぜ（笑）

## サードパーティの OpenPGP 公開鍵のインポート

今の “[Install Docker Engine on Ubuntu]” ページでは OpenPGP 公開鍵のインポート手順を以下のようにしている。

```
$ curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
```

一応説明すると，コマンドの前半で公開鍵をダウンロードし後半で /usr/share/keyrings/docker-archive-keyring.gpg ファイルに格納している。前半でダウンロードするファイルは ASCII Armor 形式のテキストなので後半の gpg コマンドで --dearmor つまりバイナリに変換して格納している。 apt-key add コマンドで鍵束ファイルに追加するのではなく，公開鍵データをそのまま単独ファイルとして置いているのがポイントである。

### ASCII Armor のままでおｋ

man コマンドで apt-key のマニュアルを見ると

> Make sure to use the "asc" extension for ASCII armored keys and the "gpg" extension for the binary OpenPGP format (also known as "GPG key public ring"). The binary OpenPGP format works for all apt versions, while the ASCII armored format works for apt version >= 1.4.

とあり，今のバージョンでは ASCII Armor 形式のまま格納しても問題ないようだ。さらに公開鍵ファイルの置き場所についても

> Recommended: Instead of placing keys into the /etc/apt/trusted.gpg.d directory, you can place them anywhere on your filesystem by using the Signed-By option in your sources.list and pointing to the filename of the key. See sources.list(5) for details. Since APT 2.4, /etc/apt/keyrings is provided as the recommended location for keys not managed by packages.

とあり，後述の Signed-By オプションと組み合わせて /etc/apt/keyrings/ ディレクトリに置くことが推奨されている。なので先程のコマンドは gpg コマンドを使うまでもなく

```
$ sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker-key.asc
```

で置き換えることができる（出力先ファイル名は適当）。簡単♪

### 【宣伝】 [gpgpdump] で公開鍵を事前確認する

拙作の [gpgpdump] コマンドを使って Web 上の任意の OpenPGP データを解析・表示できる。たとえば

```
$ gpgpdump fetch -u --indent 2 https://download.docker.com/linux/ubuntu/gpg
Public-Key Packet (tag 6) (525 bytes)
  Version: 4 (current)
  Public key creation time: 2017-02-22T18:36:26Z
  Public-key Algorithm: RSA (Encrypt or Sign) (pub 1)
  RSA public modulus n (4096 bits)
  RSA public encryption exponent e (17 bits)
User ID Packet (tag 13) (43 bytes)
  User ID: Docker Release (CE deb) <docker@docker.com>
Signature Packet (tag 2) (567 bytes)
  Version: 4 (current)
  Signiture Type: Positive certification of a User ID and Public-Key packet (0x13)
  Public-key Algorithm: RSA (Encrypt or Sign) (pub 1)
  Hash Algorithm: SHA2-512 (hash 10)
  Hashed Subpacket (33 bytes)
    Signature Creation Time (sub 2): 2017-02-22T19:34:24Z
    Key Flags (sub 27) (1 bytes)
      Flag: This key may be used to certify other keys.
      Flag: This key may be used to sign data.
      Flag: This key may be used to encrypt communications.
      Flag: This key may be used to encrypt storage.
      Flag: This key may be used for authentication.
    Preferred Symmetric Algorithms (sub 11) (4 bytes)
      Symmetric Algorithm: AES with 256-bit key (sym 9)
      Symmetric Algorithm: AES with 192-bit key (sym 8)
      Symmetric Algorithm: AES with 128-bit key (sym 7)
      Symmetric Algorithm: CAST5 (128 bit key, as per) (sym 3)
    Preferred Hash Algorithms (sub 21) (4 bytes)
      Hash Algorithm: SHA2-512 (hash 10)
      Hash Algorithm: SHA2-384 (hash 9)
      Hash Algorithm: SHA2-256 (hash 8)
      Hash Algorithm: SHA2-224 (hash 11)
    Preferred Compression Algorithms (sub 22) (4 bytes)
      Compression Algorithm: ZLIB <RFC1950> (comp 2)
      Compression Algorithm: BZip2 (comp 3)
      Compression Algorithm: ZIP <RFC1951> (comp 1)
      Compression Algorithm: Uncompressed (comp 0)
    Features (sub 30) (1 bytes)
      Flag: Modification Detection (packets 18 and 19)
    Key Server Preferences (sub 23) (1 bytes)
      Flag: No-modify
  Unhashed Subpacket (10 bytes)
    Issuer (sub 16): 0x8d81803c0ebfcd88
  Hash left 2 bytes
    b2 c9
  RSA signature value m^d mod n (4094 bits)
Public-Subkey Packet (tag 14) (525 bytes)
  Version: 4 (current)
  Public key creation time: 2017-02-22T18:36:26Z
  Public-key Algorithm: RSA (Encrypt or Sign) (pub 1)
  RSA public modulus n (4096 bits)
  RSA public encryption exponent e (17 bits)
Signature Packet (tag 2) (1086 bytes)
  Version: 4 (current)
  Signiture Type: Subkey Binding Signature (0x18)
  Public-key Algorithm: RSA (Encrypt or Sign) (pub 1)
  Hash Algorithm: SHA2-256 (hash 8)
  Hashed Subpacket (9 bytes)
    Signature Creation Time (sub 2): 2017-02-22T18:36:26Z
    Key Flags (sub 27) (1 bytes)
      Flag: This key may be used to sign data.
  Unhashed Subpacket (553 bytes)
    Issuer (sub 16): 0x8d81803c0ebfcd88
    Embedded Signature (sub 32) (540 bytes)
      Signature Packet (tag 2) (540 bytes)
        Version: 4 (current)
        Signiture Type: Primary Key Binding Signature (0x19)
        Public-key Algorithm: RSA (Encrypt or Sign) (pub 1)
        Hash Algorithm: SHA2-256 (hash 8)
        Hashed Subpacket (6 bytes)
          Signature Creation Time (sub 2): 2017-02-22T18:36:26Z
        Unhashed Subpacket (10 bytes)
          Issuer (sub 16): 0x7ea0a9c3f273fcd8
        Hash left 2 bytes
          d5 60
        RSA signature value m^d mod n (4095 bits)
  Hash left 2 bytes
    f2 b8
  RSA signature value m^d mod n (4095 bits)
```

という感じ。これでインポートする公開鍵の詳細情報を事前に確認できる。
さらに [gpgpdump] fetch コマンドに --raw オプションを付けることでフェッチしたデータをそのまま標準出力に出力できるので，先程の curl コマンドを使ったインポートの代わりに

```
$ sudo sh -c "gpgpdump fetch --raw https://download.docker.com/linux/ubuntu/gpg > /etc/apt/keyrings/docker-key.asc"
```

あるいは

```
$ gpgpdump fetch --raw https://download.docker.com/linux/ubuntu/gpg | sudo tee /etc/apt/keyrings/docker-key.asc > /dev/null
```

とすることもできる。以上宣伝でした（笑）

## Signed-By でリポジトリと公開鍵を紐付ける

今の “[Install Docker Engine on Ubuntu]” ページを参考に [Docker] Engine の APT リポジトリを source.list に登録してみる。インポートした OpenPGP 公開鍵ファイルを /etc/apt/keyrings/docker-key.asc とすると，リポジトリの登録は

```
$ echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker-key.asc] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
```

でいけるようだ。あるいは add-apt-repository コマンドを使うのであれば

```
$ sudo add-apt-repository "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker-key.asc] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
```

でもいけるかな。ポイントは `signed-by` オプション。これを使ってリポジトリと公開鍵を紐付けることができる。

ここまでできれば，あとは今までどおり apt update/install/upgrade でパッケージのインストールや更新ができる。

## サードパーティ・パッケージの公開鍵はユーザが管理する

なんでこんな面倒くさいことをするかというと，サードパーティ・パッケージの公開鍵に対応する秘密鍵が漏洩しても APT 公式側では対応できないので，鍵束から分離したいのだ。秘密鍵が漏洩した**管理されない鍵**が放置されると「信頼できる malware」を差し込まれる可能性が高くなる。

まぁ，でも，そこでサードパーティの鍵を排除するのではなく管理を分離するという発想が Linux ぽいよね。どこぞのA社とかG社とかM社とかのアプリケーション・ストアとは一線を画しているわけだ（笑）

その代わりサードパーティの公開鍵の管理はユーザ側の責務となる。 Revoke を含む公開鍵の管理をユーザ側できちんと行わないと，結局はリスクを抱え込むことになる。

気ィつけなはれや！

## 【おまけ】pgAdmin 4 のインストール

PostgreSQL サービス管理者の味方 [pgAdmin] を APT でインストールする場合も公開鍵のインポートを行う必要がある。 [Ubuntu] 22.04 LTS に対応する v6.9 が出てたので対応してみる。こんな感じだろうか（出力先ファイル名は適当）。

```
$ sudo curl -fsSL https://www.pgadmin.org/static/packages_pgadmin_org.pub -o /etc/apt/keyrings/pgadmin-4-key.asc
```

リポジトリの登録は

```
$ sudo sh -c 'echo "deb [signed-by=/etc/apt/keyrings/pgadmin-4-key.asc] https://ftp.postgresql.org/pub/pgadmin/pgadmin4/apt/$(lsb_release -cs) pgadmin4 main" > /etc/apt/sources.list.d/pgadmin4.list'
```

って感じ。あとはいつもどおりに update して upgrade すれば無問題。

[Debian]: https://www.debian.org/ "Debian -- The Universal Operating System"
[Ubuntu]: https://www.ubuntu.com/ "The leading operating system for PCs, IoT devices, servers and the cloud | Ubuntu"
[Docker]: https://www.docker.com/ "Empowering App Development for Developers | Docker"
[Install Docker Engine on Ubuntu]: https://docs.docker.com/engine/install/ubuntu/ "Install Docker Engine on Ubuntu | Docker Documentation"
[gpgpdump]: https://github.com/goark/gpgpdump "goark/gpgpdump: OpenPGP packet visualizer"
[pgAdmin]: https://www.pgadmin.org/ "pgAdmin - PostgreSQL Tools"
