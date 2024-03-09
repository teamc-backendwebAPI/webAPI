# 概要
フォーム画面に知りたい名称を入力すれば、レシピの一覧やその詳細を知ることができるアプリです。

レシピは、230万件を超えるレシピを提供している無料API（以下のサイト）から取得しています。
https://developer.edamam.com/edamam-docs-recipe-api

<br>

# エンドポイント
1. レシピ検索APIのエンドポイント<br>
 POST `/submit`　入力してもらったレシピ名を送信してもらい、その種類や詳細なレシピ情報を上記サイトから取得する。
2. レシピ詳細表示APIのエンドポイント<br>
 GET `/recipe/:index`　選択したレシピの詳細情報を取得し表示する。

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
| image       | string   | レシピの画像URL          | 
| source      | string   | レシピの出典             | 
| url         | string   | レシピURL                | 
| yield       | float64  | レシピの量               | 
| ingredients | array    | レシピの材料             | 
| calories    | float64  | レシピのカロリー         | 
| totalTime   | float64  | レシピの調理時間         | 

<br>

# エラーハンドリング
### ステータスコード:400
``` json
{
    "error":"Invalid index"
}
```
### ステータスコード:404
``` json
{
    "error":"Recipe not found"
}
```
<br>

# 認証
不要
<br>

# 使用例
https://github.com/teamc-backendwebAPI/webAPI/assets/121601977/006bc42b-5ba1-46a4-99aa-b742bc149fc3
<br>
動画のように、知りたいものを入力して使用してください。
<br>

# 制限事項
EDAMAMのDEVELOPER向けAPIを使用しています。フリーライセンスであり、検索できる回数に限りがあります。
