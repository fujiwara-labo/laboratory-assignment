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
    "github.com/gin-contrib/sessions"
    "github.com/gin-contrib/sessions/cookie"
)

func main() {
    control.DbInit()
    router := gin.Default()
    router.LoadHTMLGlob("views/*.html")
    // sessionを利用する
    store := cookie.NewStore([]byte("secret"))
    router.Use(sessions.Sessions("mysession", store))


    router.GET("/", func(c *gin.Context) {
        c.HTML(http.StatusOK, "index.html", gin.H{
             // htmlに渡す変数を定義
            "message": "test",
        })
    })
    
    // 学生ユーザー登録、ログイン画面
    router.GET("/login", func(c *gin.Context) {

        c.HTML(200, "login.html", gin.H{})
    })
    // 学生ユーザー登録
    router.POST("/signup", func(c *gin.Context) {
        var form models.Student
        // バリデーション処理
        if err := c.Bind(&form); err != nil {
            c.HTML(http.StatusBadRequest, "login.html", gin.H{"err": err})
            c.Abort()
        } else {
            student_id := c.PostForm("student_id")
            password := c.PostForm("password")
            department := c.PostForm("department")
            // 登録ユーザーが重複していた場合にはじく処理
            if err := control.CreateStudent(student_id, password, department); err != nil {
                c.HTML(http.StatusBadRequest, "login.html", gin.H{"err": err})
            }
            c.Redirect(302, "/")
        }
    })

    // 学生ユーザーログイン
    router.POST("/login", func(c *gin.Context) {
        // sessionを作成
        session := sessions.Default(c)
        session.Set("loginUser", c.PostForm("student_id"))
        session.Save()
        log.Println(session.Get("loginUser"))

        // ログインしているStudentの取得
        student := control.GetStudent(c.PostForm("student_id"))
        // DBから取得したユーザーパスワード(Hash)
        dbPassword := student.Password
        // フォームから取得したユーザーパスワード
        formPassword := c.PostForm("password")

        // ユーザーパスワードの比較
        if err := crypto.CompareHashAndPassword(dbPassword, formPassword); err != nil {
            log.Println("Could not log in")
            c.HTML(http.StatusBadRequest, "login.html", gin.H{"err": err})
            c.Abort()
        } else {
            log.Println("Could log in")
            labs := control.GetAllLab(student.Department)
            log.Println(labs)
            log.Println("collect get labs")

            c.Redirect(302, "/home-student")
        }
    })
    // 学生ユーザーログアウト（sessionクリア）
    router.POST("/logout",func(c *gin.Context) {
        session := sessions.Default(c)
        session.Clear()
        session.Save()
        log.Println(session.Get("loginUser"))
        c.Redirect(302, "/")
    })
    // 学生ユーザーホーム画面
    router.GET("/home-student", func(c *gin.Context) {
        session := sessions.Default(c)
        // if !session {
        //     c.Redirect(302, "/login")
        // }
        session_id := session.Get("loginUser")
        student_id := session_id.(string)
        c.HTML(200, "home-student.html", gin.H{
            "student_id": student_id,
        })
    })
    // 学生ユーザーの志望書提出フォーム画面
    router.GET("/form", func(c *gin.Context) {
        session := sessions.Default(c)

        session_id := session.Get("loginUser")
        student_id := session_id.(string)
        student := control.GetStudent(student_id)
        labs := control.GetAllLab(student.Department)
        c.HTML(200, "form.html", gin.H{
            "student_id": student.Student_id,
            "department": student.Department,
            "labs": labs,
        })
    })
    // フォームの取得
    router.POST("/form", func(c *gin.Context) {
        session := sessions.Default(c)

        session_id := session.Get("loginUser")
        student_id := session_id.(string)
        log.Println(student_id)
        reason := c.PostForm("reason")
        rank := c.PostForm("rank")
        lab_id := c.PostForm("lab_id")
        log.Println(lab_id)
        control.CreateAspire(student_id, lab_id, reason, rank)
        c.Redirect(302, "/home-student")
    })

    // 教員ユーザー登録、ログイン画面
    router.GET("/login-lab", func(c *gin.Context) {

        c.HTML(200, "login-lab.html", gin.H{})
    })
    // 教員ユーザー登録
    router.POST("/signup-lab", func(c *gin.Context) {
        var form models.Lab
        // バリデーション処理
        if err := c.Bind(&form); err != nil {
            c.Redirect(302, "/login-lab")
            c.Abort()
        } else {
            lab_id := c.PostForm("lab_id")
            password := c.PostForm("password")
            department := c.PostForm("department")
            // 登録ユーザーが重複していた場合にはじく処理
            if err := control.CreateLab(lab_id, password, department); err != nil {
                c.Redirect(302, "/login-lab")
            }
            c.Redirect(302, "/login-lab")
        }
    })
    router.GET("/home-lab", func(c *gin.Context) {
        session := sessions.Default(c)
        // if !session {
        //     c.Redirect(302, "/login-lab")
        // }
        session_id := session.Get("loginUser")
        lab_id := session_id.(string)
        aspires := control.GetAllAspire(lab_id)
        c.HTML(200, "home-lab.html", gin.H{
            "lab_id": lab_id,
            "lab_id2": lab_id,
            "aspires": aspires,
        })
    })
    // 教員ユーザーログイン
    router.POST("/login-lab", func(c *gin.Context) {
        // sessionを作成
        session := sessions.Default(c)
        session.Set("loginUser", c.PostForm("lab_id"))
        session.Save()
        // ログインしているStudentの取得
        lab := control.GetLab(c.PostForm("lab_id"))
        log.Println(lab)
        // DBから取得したユーザーパスワード(Hash)
        dbPassword := lab.Password
        log.Println(dbPassword)
        // フォームから取得したユーザーパスワード
        formPassword := c.PostForm("password")

        // ユーザーパスワードの比較
        if err := crypto.CompareHashAndPassword(dbPassword, formPassword); err != nil {
            log.Println("Could not log in")
            c.Redirect(302, "/login-lab")
            c.Abort()
        } else {
            log.Println("Could log in")
            c.Redirect(302, "/home-lab")
        }
    })

    router.Run(":8080")

}