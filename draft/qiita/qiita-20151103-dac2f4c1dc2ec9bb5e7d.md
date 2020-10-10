---
title: "NYAGOS で Lua"
emoji: "😀"
type: "tech"
topics: [Windows,Lua,nyagos]
published: false
---
この記事は [text.Baldanders.info の記事](http://text.baldanders.info/remark/2015/nyagos-and-lua/)を再掲載したものです。[Lua] 初心者なので不備等あると思いますがご容赦を。

[NYAGOS] の特徴のひとつに， [Lua] のエンジンを内蔵し解釈できる点がある。というか [NYAGOS] 自身が機能拡張を [Lua] で実装していて，ユーザも hackable に [NYAGOS] をカスタマイズできる。

また `nyagos.exe` を `-f` オプションを付けて起動すると [Lua] のソースファイルをスクリプトとして実行できる。たとえば以下のコード `hello.lua` を以下のように記述し

```lua
print "Hello World!"
```

`nyagos.exe` で実行すると

```
C:>nyagos.exe -f hello.lua
Hello World!
```

となる。

これだけだと大したことはないが， `nyagos.exe` でスクリプトを実行した場合は [NYAGOS] の [Lua] 拡張が使える。例えばカレント・フォルダを取得する `getwd.lua` を以下のように記述し

```lua
nyagos.write(nyagos.getwd().."\n")
```

`nyagos.exe` で実行すると

```
C:\Users\username>nyagos.exe -f getwd.lua
C:\Users\username
```

となる。

[Lua] 拡張については [NYAGOS のドキュメント](https://github.com/zetamatta/nyagos/blob/master/Doc/nyagos_ja.md)を参照するとよいだろう。特筆すべきは `alias` 機能で，単純なコマンドの置き換えだけでなく

```
nyagos.alias.cmdname = function(args) ... end
```

のような形式で新しい内部コマンドを作成することもできる。

ところで `nyagos.exe` を普通に起動して `ls -oF` コマンドでファイルの一覧を表示させてみると，拡張子 `.lua` のファイルは実行可能ファイルになっていることが分かる。ただ，実際に起動しようとすると

```
~> hello.lua
'lua' is not recognized as an internal or external command,
operable program or batch file
```

と怒られた。どうやら `lua` コマンドがねーよ，と言っているらしい（確かに [Lua] の実行モジュールは入れてないのだが）。

[NYAGOS] をインストールしたフォルダにある `nyagos.d\suffix.lua` ファイルの末尾あたりを見ると

```lua
suffix.pl="perl"
if nyagos.which("ipy") then
  suffix.py="ipy"
elseif nyagos.which("py") then
  suffix.py="py"
else
  suffix.py="python"
end
suffix.rb="ruby"
suffix.lua="lua"
suffix.awk={"awk","-f"}
suffix.js={"cscript","//nologo"}
suffix.vbs={"cscript","//nologo"}
suffix.ps1={"powershell","-file"}
```

となっていて，拡張子 `.lua` のファイルが `lua` コマンドに関連付けられているのが分かる。ってことは，これを `nyagos.exe` に書き換えればいいわけだ。とはいえ `nyagos.d\suffix.lua` ファイルを直接いじるわけにはいかないので（バージョンアップのたびに上書きされる）， `%HOME%` または `%HOMEPATH%` フォルダにある `.nyagos` ファイルに以下の記述を追加する[^a]。

 ```lua
 suffix.lua={"nyagos.exe","-f"}
 ```

 [^a]: `suffix` ではなく `alias` で `lua` コマンドを定義する手もある。

これで `nyagos.exe` を起動し直して `hello.lua` ファイルを実行してみる。

```
~> hello.lua
Hello World!
```

おおっ，うまくいった。これで [Lua] スクリプトを [NYAGOS] 上でバッチファイルのように扱うことができる。

さぁて，これでいよいよもって [Lua] の勉強をしないといけなくなった。Pascal 系の構文は目が滑るんだよなぁ。

## ブックマーク

### NYAGOS 作者による解説

- [環境変数PATHが長すぎて、追加したパスが有効にならないぜ！ - Qiita](http://qiita.com/zetamatta/items/a49e3a40201511128508)
- [nyagosスクリプト解説 - svn のサブコマンドを勝手に拡張する - Qiita](http://qiita.com/zetamatta/items/c4ad3cc55c5afa74da63)
- [nyagosスクリプト解説 - ezoe.lua「コマンドではない。」 - Qiita](http://qiita.com/zetamatta/items/29a85695813926cafd2c)
- [nyagosスクリプト解説 - 逆クォートによるコマンド出力展開編(backquote.lua) - Qiita](http://qiita.com/zetamatta/items/cdff310f53faf3369e48)
- [nyagosスクリプト解説 - CMD.EXEで化けさせず、nyagosの中だけプロンプトをカラー化 - Qiita](http://qiita.com/zetamatta/items/c08586c85fa73c182a7a)
- [nyagosスクリプト解説 - VisualStudio れんけー - Qiita](http://qiita.com/zetamatta/items/89a907f4bd46d1750c31)
- [git 付属のPerlやunzipとかを使いたいけどsort,find,lnとかは要らない - Qiita](http://qiita.com/zetamatta/items/1fe83f736b0254e02415) [^b]

[^b]: 最近の [Git for Windows](https://git-for-windows.github.io/) は bash 関連のコマンドを `Git\usr\bin` フォルダに集めている。 git コマンドだけが必要なら `Git\cmd` フォルダにのみパスを通せばよい。 bash を使う場合は `Git\bin` フォルダにある `bash.exe` を起動するのが一番安全なようだ。

### Lua の解説

- [Luaプログラミング入門 | densan-labs.net](http://densan-labs.net/tech/lua/index.html)
- [高速スクリプト言語 [Lua] を始めよう — WTOPIA v1.0 documentation](http://www.ie.u-ryukyu.ac.jp/~e085739/lua.hajime.html)
- [良いもの。悪いもの。: Lua基礎文法最速マスター](http://handasse.blogspot.com/2010/02/lua.html)
- [紀子さん＠へぼぷろぐらまの日常 | Luaで日付時間操作。](http://noriko3.blog42.fc2.com/blog-entry-128.html)
- [Luaでファイルの読み書きを行なう - Symfoware](http://symfoware.blog68.fc2.com/blog-entry-454.html)

[ConEmu]: https://conemu.github.io/ "ConEmu - Handy Windows Terminal"
[NYAGOS]: http://www.nyaos.org/index.cgi?p=NYAGOS "NYAOS.ORG - NYAGOS"
[Lua]: http://www.lua.org/ "The Programming Language Lua"

