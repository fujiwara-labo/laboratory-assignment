package main

import (
	"log"
	"net/http"
    "github.com/fujiwara-labo/laboratory-assignment.git/crypto"
    "github.com/fujiwara-labo/laboratory-assignment.git/models"
    "github.com/fujiwara-labo/laboratory-assignment.git/control"
    _ "github.com/go-sql-driver/mysql"

    "github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
    control.DbInit()
    router := gin.Default()
	router.LoadHTMLGlob("views/*.html")

    router.GET("/", func(c *gin.Context) {
        c.HTML(http.StatusOK, "index.html", gin.H{
             // htmlに渡す変数を定義
            "message": "test",
        })
    })
    
    // ユーザー登録、ログイン画面
    router.GET("/login", func(c *gin.Context) {

        c.HTML(200, "login.html", gin.H{})
    })
    // ユーザー登録
    router.POST("/signup", func(c *gin.Context) {
        var form models.User
        // バリデーション処理
        if err := c.Bind(&form); err != nil {
            c.HTML(http.StatusBadRequest, "login.html", gin.H{"err": err})
            c.Abort()
        } else {
            username := c.PostForm("username")
            password := c.PostForm("password")
            // 登録ユーザーが重複していた場合にはじく処理
            if err := control.CreateUser(username, password); err != nil {
                c.HTML(http.StatusBadRequest, "signup.html", gin.H{"err": err})
            }
            c.Redirect(302, "/")
        }
    })

    // ユーザーログイン
    router.POST("/login", func(c *gin.Context) {

        // DBから取得したユーザーパスワード(Hash)
        dbPassword := control.GetUser(c.PostForm("username")).Password
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
            c.HTML(200, "home-student.html", gin.H{})
        }
    })

    router.Run(":8080")

}