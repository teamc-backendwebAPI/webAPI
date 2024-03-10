# 概要

フォーム画面に知りたい名称を入力すれば、レシピの一覧やその詳細を知ることができるアプリです。<br>

レシピは、230 万件を超えるレシピを提供している無料 API（以下のサイト）から取得しています。
https://developer.edamam.com/edamam-docs-recipe-api
<br>

# 環境構築(MySQL)

mysql を起動してターミナルで以下のコマンドを実行してデータベースをローカル環境に作成してください。


1. データベースの作成
   `CREATE DATABASE c_auth;`

2. 新しいユーザーの作成
   `CREATE USER 'c_auth'@'localhost' IDENTIFIED BY 'c_webAPI_@auth123';`

3. 権限の付与
   `GRANT ALL PRIVILEGES ON c_auth.* TO 'c_auth'@'localhost';`

4. 権限のリフレッシュ
   `FLUSH PRIVILEGES;`

.envファイルをルートディレクトリに作って以下のように記述してください。""はいらないです。
```
DB_DBMS=mysql
DB_USER=c_auth
DB_PASS=c_webAPI_@auth123
DB_DBNAME=c_auth
```

`go get`コマンドで不足しているライブラリはインポートしてから`go run .`をmain.goがあるディレクトリで実行してください。

## windowsの場合
   `mysql --user=root --password`
   で接続可能
   

# エンドポイント

1. レシピ検索 API のエンドポイント<br>
   POST `/submit`　入力してもらったレシピ名を送信してもらい、その種類や詳細なレシピ情報を上記サイトから取得する。
2. レシピ詳細表示 API のエンドポイント<br>
   GET `/recipe/:index`　選択したレシピの詳細情報を取得し表示する。
3. APIドキュメントの表示

4. レシピデータのページ分け

5. レシピデータの詳細表示

<br>

# リクエスト

### リクエストボディ

| パラメータ | データ型 | 説明                 |
| ---------- | -------- | -------------------- |
| name       | string   | 検索するレシピの名前 |

# レスポンス

### ステータスコード

| ステータスコード | 説明                 |
| ---------------- | -------------------- |
| 200              | リクエスト成功       |
| 400              | 無効なインデックス   |
| 404              | レシピが見つからない |

### レスポンスボディ

| パラメータ  | データ型 | 説明                     |
| ----------- | -------- | ------------------------ |
| recipes     | array    | 検索したレシピを含む配列 |
| recipe      | object   | レシピの情報             |
| label       | string   | レシピ名                 |
| image       | string   | レシピの画像 URL         |
| source      | string   | レシピの出典             |
| url         | string   | レシピ URL               |
| yield       | float64  | レシピの量               |
| ingredients | array    | レシピの材料             |
| calories    | float64  | レシピのカロリー         |
| totalTime   | float64  | レシピの調理時間         |

<br>

# エラーハンドリング

### ステータスコード:400

```json
{
  "error": "Invalid index"
}
```

### ステータスコード:404

```json
{
  "error": "Recipe not found"
}
```

<br>

# 認証

不要
<br>
# こだわったところ

# 改善点


# 制限事項

EDAMAM の DEVELOPER 向け API を使用しています。フリーライセンスであり、検索できる回数に限りがあります。
