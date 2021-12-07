package main

import (
	"log"
	"net/http"
    "os"
    "github.com/fujiwara-labo/laboratory-assignment.git/crypto"
    "github.com/fujiwara-labo/laboratory-assignment.git/models"
    "github.com/fujiwara-labo/laboratory-assignment.git/control"
    _ "github.com/go-sql-driver/mysql"

    "github.com/gin-gonic/gin"
    "github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
)

// User モデルの宣言
// type User struct {
// 	gorm.Model
// 	Username string `form:"username" binding:"required" gorm:"unique;not null"`
// 	Password string `form:"password" binding:"required"`
// }

// func gormConnect() *gorm.DB {
// 	err := godotenv.Load()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
//     DBMS := os.Getenv("DRIVER")
//     CONNECT := os.Getenv("DSN")
// 	// USER := os.Getenv("MYSQL_USER")
// 	// PASS := os.Getenv("MYSQL_PASSWORD")
// 	// DBNAME := os.Getenv("MYSQL_DATABASE")
// 	// CONNECT := USER + ":" + PASS + "@/" + DBNAME + "?parseTime=true"
// 	db, err := gorm.Open(DBMS, CONNECT)
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	return db
// }
// // DBの初期化
// func dbInit() {
// 	db := gormConnect()
// 	// コネクション解放
// 	defer db.Close()
//     //構造体に基づいてテーブルを作成
// 	db.AutoMigrate(&User{})
// }

// ユーザー登録処理
// func createUser(username string, password string) []error {
// 	passwordEncrypt, _ := crypto.PasswordEncrypt(password)
// 	db := gormConnect()
// 	defer db.Close()
// 	// Insert処理
// 	if err := db.Create(&User{Username: username, Password: passwordEncrypt}).GetErrors(); err != nil {
// 		return err
// 	}
// 	return nil

// }
// // ユーザーを一件取得
// func getUser(username string) User {
// 	db := gormConnect()
// 	var user User
// 	db.First(&user, "username = ?", username)
// 	db.Close()
// 	return user
// }

func main() {
    router := gin.Default()
	router.LoadHTMLGlob("views/*.html")

    control.dbInit()

    router.GET("/", func(c *gin.Context) {
        c.HTML(http.StatusOK, "index.html", gin.H{
             // htmlに渡す変数を定義
            "message": "test",
        })
    })
    
    // ユーザー登録画面
    router.GET("/signup", func(c *gin.Context) {

        c.HTML(200, "signup.html", gin.H{})
    })

    // ユーザー登録
    router.POST("/signup", func(c *gin.Context) {
        var form models.User
        // バリデーション処理
        if err := c.Bind(&form); err != nil {
            c.HTML(http.StatusBadRequest, "signup.html", gin.H{"err": err})
            c.Abort()
        } else {
            username := c.PostForm("username")
            password := c.PostForm("password")
            // 登録ユーザーが重複していた場合にはじく処理
            if err := control.createUser(username, password); err != nil {
                c.HTML(http.StatusBadRequest, "signup.html", gin.H{"err": err})
            }
            c.Redirect(302, "/")
        }
    })

    // ユーザーログイン画面
    router.GET("/login", func(c *gin.Context) {

        c.HTML(200, "login.html", gin.H{})
    })

    // ユーザーログイン
    router.POST("/login", func(c *gin.Context) {

        // DBから取得したユーザーパスワード(Hash)
        dbPassword := control.getUser(c.PostForm("username")).Password
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
    })

    router.Run(":8080")

}