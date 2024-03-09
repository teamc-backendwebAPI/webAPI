package main

import (
    "database/sql"
    "fmt"
    "log"
    "net/http"

    _ "github.com/go-sql-driver/mysql"
    "golang.org/x/crypto/bcrypt"
)

var db *sql.DB
var err error

func main() {
    // データベースへの接続
    db, err = sql.Open("mysql", "ユーザー名:パスワード@tcp(127.0.0.1:3306)/login_system")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // ルート定義
    http.HandleFunc("/register", registerHandler)
    http.HandleFunc("/login", loginHandler)

    // サーバー起動
    fmt.Println("Starting server on :8080...")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" {
        http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
        return
    }

    username := r.FormValue("username")
    password := r.FormValue("password")

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        log.Fatal(err)
    }

    _, err = db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", username, string(hashedPassword))
    if err != nil {
        http.Error(w, "Failed to register user", http.StatusInternalServerError)
        return
    }

    fmt.Fprintf(w, "User registered successfully")
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" {
        http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
        return
    }

    username := r.FormValue("username")
    password := r.FormValue("password")

    var hashedPassword string
    err = db.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&hashedPassword)
    if err != nil {
        if err == sql.ErrNoRows {
            http.Error(w, "Username or password is incorrect", http.StatusUnauthorized)
            return
        }
        log.Fatal(err)
    }

    err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
    if err != nil {
        http.Error(w, "Username or password is incorrect", http.StatusUnauthorized)
        return
    }

    fmt.Fprintf(w, "Logged in successfully")
}
