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
    
    // 学生ユーザー登録、ログイン画面
    router.GET("/admin", func(c *gin.Context) {

        c.HTML(200, "admin.html", gin.H{})
    })
    // 学生ユーザー登録
    router.POST("/signup", func(c *gin.Context) {
        var form models.Student
        // バリデーション処理
        if err := c.Bind(&form); err != nil {
            c.HTML(http.StatusBadRequest, "admin.html", gin.H{"err": err})
            c.Abort()
        } else {
            student_id := c.PostForm("student_id")
            password := c.PostForm("password")
            department := c.PostForm("department")
            // 登録ユーザーが重複していた場合にはじく処理
            if err := control.CreateStudent(student_id, password, department); err != nil {
                c.HTML(http.StatusBadRequest, "admin.html", gin.H{"err": err})
            }
            c.Redirect(302, "/")
        }
    })

    // 学生ユーザーログイン
    router.POST("/login", func(c *gin.Context) {

        // DBから取得したユーザーパスワード(Hash)
        student := control.GetStudent(c.PostForm("student_id"))
        dbPassword := student.Password
        // フォームから取得したユーザーパスワード
        formPassword := c.PostForm("password")

        // ユーザーパスワードの比較
        if err := crypto.CompareHashAndPassword(dbPassword, formPassword); err != nil {
            log.Println("Could not log in")
            c.HTML(http.StatusBadRequest, "admin.html", gin.H{"err": err})
            c.Abort()
        } else {
            log.Println("Could log in")
            labs := control.GetAllLab(student.Department)
            log.Println(labs)
            log.Println("collect get labs")

            c.HTML(200, "home-student.html", gin.H{
                "student_id": student.Student_id,
                "department": student.Department,
                "labs": labs,
            })
        }
    })

    // 教員ユーザー登録、ログイン画面
    router.GET("/admin-lab", func(c *gin.Context) {

        c.HTML(200, "admin-lab.html", gin.H{})
    })
    // 教員ユーザー登録
    router.POST("/signup-lab", func(c *gin.Context) {
        var form models.Lab
        // バリデーション処理
        if err := c.Bind(&form); err != nil {
            c.HTML(http.StatusBadRequest, "admin-lab.html", gin.H{"err": err})
            c.Abort()
        } else {
            lab_id := c.PostForm("lab_id")
            password := c.PostForm("password")
            department := c.PostForm("department")
            // 登録ユーザーが重複していた場合にはじく処理
            if err := control.CreateLab(lab_id, password, department); err != nil {
                c.HTML(http.StatusBadRequest, "admin-lab.html", gin.H{"err": err})
            }
            c.Redirect(302, "/")
        }
    })

    // 教員ユーザーログイン
    router.POST("/login-lab", func(c *gin.Context) {

        // DBから取得したユーザーパスワード(Hash)
        dbPassword := control.GetLab(c.PostForm("lab_id")).Password
        log.Println(dbPassword)
        // フォームから取得したユーザーパスワード
        formPassword := c.PostForm("password")

        // ユーザーパスワードの比較
        if err := crypto.CompareHashAndPassword(dbPassword, formPassword); err != nil {
            log.Println("Could not log in")
            c.HTML(http.StatusBadRequest, "admin-lab.html", gin.H{"err": err})
            c.Abort()
        } else {
            log.Println("Could log in")
            c.HTML(200, "home-lab.html", gin.H{})
        }
    })

    router.Run(":8080")

}