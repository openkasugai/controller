# OpenKasugai Controller へのコントリビューション

本コントローラは、CNCF等の関連する既存OSSプロジェクトへの[DCIコンセプト](https://www.rd.ntt/research/JN202311_23718.html)の提案に向けたリファレンスとして公開しています。今後は、既存OSSプロジェクト内での議論・コントリビューションを予定しています。

一方で、もしOpenKasugai Controllerに対するコントリビューションの希望が特にありましたら GitHubの[Discussions](https://github.com/openkasugai/controller/discussions)機能を使用し、ご提案ください。

## バグの報告・機能追加要望

GitHub の [Issues](https://github.com/openkasugai/controller/issues) を使用してください。

### バグ

- Labels には `bug` を付与してください。
- 先頭のコメントに以下を記載してください。
  - 環境
    - OpenKasugai ControllerのVersion
    - Node構成、OS Kernel version, 使用デバイスの型番など
  - 問題事象
    - 具体的な問題発生状況の説明
  - 再現手順 (optional)
    - 再現手順が複雑な場合は記載してください。

### 機能追加

- Labels には `enhancement` を付与してください。
- 先頭のコメントに以下を記載してください。
  - 機能仕様
    - 実装する最終仕様を記載してください。
      - 議論により仕様を変更した場合は、先頭のコメントに最終仕様が確認できるようにしてください。
  - 参考
    - 議論内容の分かるDiscussionsなどのリンクを記載してください。

### ドキュメント改善

- Labels には `documentation` を付与してください。
- 先頭のコメントに以下を記載してください。
  - 概要
    - 変更対象文書と修正内容を記載してください。

## その他質問、設計に関する要望やアイデアなど

GitHub の [Discussions](https://github.com/openkasugai/controller/discussions)機能をご利用ください。

## プルリクエスト

- コントリビュータは、最初に OpenKasugai Controller リポジトリの main branch を fork してください。
- コントリビュータは、fork したリポジトリ上で topic branch を作成し、OpenKasugai Controller 配下のリポジトリの main branch に対して Pull Request を行ってください。
  - topic branch のブランチ名は任意です。
- コントリビュータは、[DCO](https://developercertificate.org/)に同意する必要があります。
  - DCO に同意していることを示すため、全てのコミットに対して、コミットメッセージに以下を記入してください。
    - `Signed-off-by: Random J Developer <random@developer.example.org>`
      - 氏名の部分は、本名を使用してください。
      - GitHub の Profile の Name に同じ名前、Email に同じメールアドレスを設定する必要があります。
      - `git commit -s` でコミットに署名を追加ください。
- Pull Request を発行する際は、対応する Issue に紐づけてください。
  - 対応する Issue がない場合は Pull Request の発行前に作成してください。
- Pull Request のタイトルには、"fix"に続いて対処した issue 番号および修正の概要を記入してください。
  - `fix #[issue番号] [修正の概要]`
- Pull Request の本文は、テンプレートを使用してください。
- Pull Requestのマージにはスタイルチェック(GitHub Actions: golangci-lint)の成功、およびメンテナ2名以上のApproveが必要です。

## リリースサイクル

新規機能追加やバグフィックスなどの変更があれば、毎年10月、4月にリリースを行います。
変更がなければリリースを行いません。
