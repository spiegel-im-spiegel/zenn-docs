---
title: "私の欲しいエラーと貴方の欲しいエラーは違う"
---

エラーハンドリングで最も考慮すべきことは

💡 **利用者が欲しいエラー情報と提供者が欲しいエラー情報は異なる** 💡

という点だろう。

エラーが発生した際に利用者が最も欲しい情報は「どうすればいいのか？」である。そのためのヒントとして「何故エラーが起こったのか？」も欲しいわけだ。

たとえば，コマンドライン・ツールのフレームワークを提供する [spf13/cobra] は利用者がコマンド入力を間違えた際に正しいコマンドを推測して教えてくれる。[私が公開しているコマンドライン・ツール](https://github.com/spiegel-im-spiegel/gpgpdump "spiegel-im-spiegel/gpgpdump: OpenPGP packet visualizer")だとこんな感じ。

```
$ gpgpdump http
Error: unknown command "http" for "gpgpdump"

Did you mean this?
	hkp

Run 'gpgpdump --help' for usage.
```

これで本当に `gpgpdump hkp` と打ち間違えたのなら打ち直せばいいし，全然違うというのなら `gpgpdump --help` で使い方を表示してみればいい，と分かる。

一方，提供者側にとって（利用者にエラー情報提供するにせよ）最も欲しい情報は「どうやって起こったか？」である。これを知るためには，「何故？」よりむしろ，エラーが発生時の「文脈」をできるだけかき集めることが重要である。

エラー発生時にスタック情報を欲しがるエンジニアが多いのは，この情報が「文脈」の一部となりうるからだ。でも，これは個人的な見解だが，スタック情報は9割以上がノイズである（実行中のプログラムの構造解析がしたいなら別だが）。喩えるなら藁束の中から金の針を探すようなものだ。

じゃあ，エラーハンドリングはどういう戦術をとるのがいいのだろう。

[Go]: https://golang.org/ "The Go Programming Language"
[spf13/cobra]: https://github.com/spf13/cobra "spf13/cobra: A Commander for modern Go CLI interactions"
<!-- eof -->
