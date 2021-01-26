# おそるべき絵文字

Zenn 記事投稿時にちょっと気になることがあって issue を投げたのだが

https://github.com/zenn-dev/zenn-roadmap/issues/224

じゃあ，どこからどこまでが絵文字なんだろうと気になってちょっと調べてみた。

**情報歓迎**


絵文字 #️⃣ は

| Unicode Point | 字形 | Unicode 名称               |
| ------------- |:----:| -------------------------- |
| U+0023        | `#`  | NUMBER SIGN                |
| U+FE0F        |      | VARIATION SELECTOR-16      |
| U+20E3        | ` ⃣` | COMBINING ENCLOSING KEYCAP |

という合成列らしい。えー。絵文字異体字セレクタの後ろに更に結合文字が付くのかよ。しかも既定文字はただの半角記号だよ orz

これはヤバい匂いしかしない。興味本位で手を出していい案件じゃなかった！！

おそるべき絵文字の闇を垣間見てしまったわけだが，これで思い出した記事がある。

https://text.baldanders.info/remark/2017/12/character-of-the-new-era-name/

まぁ「技術的負債」と口走ったのは[私の黒歴史](https://text.baldanders.info/remark/2020/12/technical-debt-and-hacker/ "技術的負債とハッカー")としてスルーしていただけるとありがたいが，
よく考えたら「㍻ U+337B」とかを「文字」だと思うから不合理に感じるのであって
「絵文字」と思えばアリなのか。

ちなみに，先日本家ブログで書いた

https://text.baldanders.info/golang/unicode-rangetables/

で確認してみたら「㍻ U+337B」は "Synbol/other" に分類されるようだ。つまり（絵文字ではないが）絵文字と同等ということ。
