---
title: "はじめに"
---

## はじめに

少し前に [PostgreSQL] サービスに [Go] でアクセスする方法についてちょっとした調べものをした。そのときの作業メモをブログ記事として残そうと思ったのだが，単ページで収まりそうになかったので Zenn 本の体裁で書き記しておく。体裁は「本」だが，中身はただの作業記録である。ちゃんとした解説をご所望の方にはあしからずご了承のほどを。講釈はいいから動くコードをくれ！ という方には多少なりと参考になるかもしれない。

なお，筆者の作業環境は基本的に Linux/Ubuntu なので例示についてもこれを前提に記述するが，他プラットフォームの方は適当に読み替えていただけるとありがたい。また，この本に出てくるサンプルコードは以下のリポジトリで公開している。

https://github.com/spiegel-im-spiegel/ent-sample
https://github.com/spiegel-im-spiegel/gorm-sample

## ElephantSQL の利用

今回 [PostgreSQL] について調べるにあたって（仕事に絡まない）手頃な PaaS がないかとググってみたら [ElephantSQL] というのがよさそうである。

https://www.elephantsql.com/

きちんと開発環境を整えるなら Docker という手もあるし，本格的な運用をするなら大手クラウド・プロバイダ企業の PaaS を使うほうが何かと手厚いのかもしれないが，ちょろんと調べるだけなら [ElephantSQL] の無料プランで十分賄えるし，不特定相手ではない小規模運用なら有料プランも考慮の余地があるかもしれない。 [ElephantSQL] の利用方法については以下の Qiita 記事が参考になった。ありがとう！

https://qiita.com/mikankitten/items/a9a0363c7b455e928179

なお，2021年9月時点では AWS の東京リージョンも利用できるようだ。この本では [ElephantSQL] のサービスが利用可能であることを前提に話を進めていく。

[Go]: https://go.dev/
[PostgreSQL]: https://www.postgresql.org/ "PostgreSQL: The world's most advanced open source database"
[ElephantSQL]: https://www.elephantsql.com/ "ElephantSQL - PostgreSQL as a Service"
