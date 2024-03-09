package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
    // .envファイルから環境変数を読み込む
    err := godotenv.Load("../../.env")
    if err != nil {
        log.Fatal("Error loading .env file")
    }
    // 環境変数からデータベース接続設定を取得
    dbUser := os.Getenv("DB_USER")
    dbPassword := os.Getenv("DB_PASSWORD")
    dbHost := os.Getenv("DB_HOST")
    dbPort := os.Getenv("DB_PORT")
    dbName := os.Getenv("DB_NAME")

    // データベース接続設定
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // レシピの詳細を取得
    var (
        name        string
        description string
        difficulty  string
        prepTime    string
        cookTime    string
    )
    recipeName := "チキンカレー"
    err = db.QueryRow("SELECT name, description, difficulty, prep_time, cook_time FROM recipes WHERE name = ?", recipeName).Scan(&name, &description, &difficulty, &prepTime, &cookTime)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("レシピ名: %s\n説明: %s\n難易度: %s\n準備時間: %s\n調理時間: %s\n", name, description, difficulty, prepTime, cookTime)

    // 材料リストを取得
    rows, err := db.Query(`
        SELECT i.name, ri.amount
        FROM ingredients i
        JOIN recipe_ingredients ri ON i.id = ri.ingredient_id
        JOIN recipes r ON ri.recipe_id = r.id
        WHERE r.name = ?`, recipeName)
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()
    fmt.Println("材料リスト:")
    for rows.Next() {
        var ingredient, amount string
        if err := rows.Scan(&ingredient, &amount); err != nil {
            log.Fatal(err)
        }
        fmt.Printf("- %s: %s\n", ingredient, amount)
    }
    if err := rows.Err(); err != nil {
        log.Fatal(err)
    }

    // 調理手順を取得
    steps, err := db.Query(`
        SELECT step_order, description
        FROM recipe_steps
        JOIN recipes ON recipe_steps.recipe_id = recipes.id
        WHERE recipes.name = ?
        ORDER BY step_order`, recipeName)
    if err != nil {
        log.Fatal(err)
    }
    defer steps.Close()
    fmt.Println("調理手順:")
    for steps.Next() {
        var stepOrder int
        var description string
        if err := steps.Scan(&stepOrder, &description); err != nil {
            log.Fatal(err)
        }
        fmt.Printf("手順 %d: %s\n", stepOrder, description)
    }
    if err := steps.Err(); err != nil {
        log.Fatal(err)
    }
}
