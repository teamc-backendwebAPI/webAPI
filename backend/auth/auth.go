package auth

import (
	// "log"

	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql" //直接的な記述が無いが、インポートしたいものに対しては"_"を頭につける決まり
	"github.com/jinzhu/gorm"

	"github.com/joho/godotenv"

	"github.com/ymdd1/mytweet/crypto"
)

// User モデルの宣言
type User struct {
	gorm.Model
	Username string `form:"username" binding:"required" gorm:"unique;not null"`
	Password string `form:"password" binding:"required"`
}

func gormConnect() *gorm.DB {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal(err)
	}
	DBMS := os.Getenv("DB_DBMS")
	USER := os.Getenv("DB_USER")
	PASS := os.Getenv("DB_PASS")
	DBNAME := os.Getenv("DB_DBNAME")
	CONNECT := USER + ":" + PASS + "@/" + DBNAME + "?parseTime=true"
	db, err := gorm.Open(DBMS, CONNECT)
	if err != nil {
		panic(err.Error())
	}
	return db
}

// DBの初期化
func DbInit() {
	db := gormConnect()
	// コネクション解放
	defer db.Close()
	db.AutoMigrate(&User{})
}

// ユーザー登録処理
func createUser(username string, password string) []error {
	passwordEncrypt, _ := crypto.PasswordEncrypt(password)
	db := gormConnect()
	defer db.Close()
	// Insert処理
	if err := db.Create(&User{Username: username, Password: passwordEncrypt}).GetErrors(); err != nil {
		return err
	}
	return nil
}

// ユーザーログイン処理
func LoginUser(c *gin.Context) {
	// DBから取得したユーザーパスワード(Hash)
	dbPassword := getUser(c.PostForm("username")).Password
	log.Println(dbPassword)
	// フォームから取得したユーザーパスワード
	formPassword := c.PostForm("password")

	// ユーザーパスワードの比較
	if err := crypto.CompareHashAndPassword(dbPassword, formPassword); err != nil {
		log.Println("ログインできませんでした")
		c.HTML(http.StatusBadRequest, "login.html", gin.H{"err": err})
		c.Abort()
	} else {
		log.Println("ログインできました")
		c.Redirect(302, "/")
	}
}

// ユーザーを一件取得
func getUser(username string) User {
	db := gormConnect()
	var user User
	db.First(&user, "username = ?", username)
	db.Close()
	return user
}

func SignUpUser(c *gin.Context) {
	var form User
		// バリデーション処理
		if err := c.Bind(&form); err != nil {
			c.HTML(http.StatusBadRequest, "signup.html", gin.H{"err": err})
			c.Abort()
		} else {
			username := c.PostForm("username")
			password := c.PostForm("password")
			log.Println(username)
			// 登録ユーザーが重複していた場合にはじく処理
			if err := createUser(username, password); err != nil {
				log.Println(username)
				c.Redirect(http.StatusFound, "/login")
				c.Abort()
			}
			log.Println("ユーザー登録できました")
		}
}

func SomePageHandler(c *gin.Context) {
	// セッションやクッキーからユーザー名を取得する例
	username, exists := c.Get("username") // これは仮のコードです。実際の取得方法に置き換えてください。

	if exists {
			c.HTML(http.StatusOK, "some_template.html", gin.H{
					"username": username,
			})
	} else {
			c.HTML(http.StatusOK, "some_template.html", gin.H{
					"username": nil,
			})
	}
}


// func main() {
// 	router := gin.Default()
// 	router.LoadHTMLGlob("../../frontend/auth/*.html")

// 	dbInit()

// 	// ユーザー登録画面
// 	router.GET("/signup", func(c *gin.Context) {

// 		c.HTML(200, "signup.html", gin.H{})
// 	})

// 	// ユーザー登録
// 	router.POST("/signup", func(c *gin.Context) {
// 		var form User
// 		// バリデーション処理
// 		if err := c.Bind(&form); err != nil {
// 			c.HTML(http.StatusBadRequest, "signup.html", gin.H{"err": err})
// 			c.Abort()
// 		} else {
// 			username := c.PostForm("username")
// 			password := c.PostForm("password")
// 			// 登録ユーザーが重複していた場合にはじく処理
// 			if err := createUser(username, password); err != nil {
// 				c.HTML(http.StatusBadRequest, "signup.html", gin.H{"err": err})
// 			}
// 			c.Redirect(302, "/")
// 		}
// 	})

// 	// ユーザーログイン画面
// 	router.GET("/login", func(c *gin.Context) {

// 		c.HTML(200, "login.html", gin.H{})
// 	})

// 	// ユーザーログイン
// 	router.POST("/login", func(c *gin.Context) {

// 		// DBから取得したユーザーパスワード(Hash)
// 		dbPassword := getUser(c.PostForm("username")).Password
// 		log.Println(dbPassword)
// 		// フォームから取得したユーザーパスワード
// 		formPassword := c.PostForm("password")

// 		// ユーザーパスワードの比較
// 		if err := crypto.CompareHashAndPassword(dbPassword, formPassword); err != nil {
// 			log.Println("ログインできませんでした")
// 			c.HTML(http.StatusBadRequest, "login.html", gin.H{"err": err})
// 			c.Abort()
// 		} else {
// 			log.Println("ログインできました")
// 			c.Redirect(302, "/")
// 		}
// 	})

// 	router.Run()
// }
