package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/fujiwara-labo/laboratory-assignment.git/control"
	"github.com/fujiwara-labo/laboratory-assignment.git/crypto"
	"github.com/fujiwara-labo/laboratory-assignment.git/models"
	_ "github.com/go-sql-driver/mysql"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
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
	// 管理者ユーザー登録、ログイン画面
	router.GET("/login-admin", func(c *gin.Context) {

		c.HTML(200, "login-admin.html", gin.H{})
	})
	// 管理者ユーザーログアウト
	router.POST("/logout-admin", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Clear()
		session.Save()
		log.Println(session.Get("loginUser"))
		c.Redirect(302, "/login-admin")
	})
	// 管理者ユーザー登録
	router.POST("/signup-admin", func(c *gin.Context) {
		var form models.Admin
		// バリデーション処理
		if err := c.Bind(&form); err != nil {
			c.HTML(http.StatusBadRequest, "login-admin.html", gin.H{"err": err})
			c.Abort()
		} else {
			admin_id := c.PostForm("admin_id")
			password := c.PostForm("password")
			// 登録ユーザーが重複していた場合にはじく処理
			if err := control.CreateAdmin(admin_id, password); err != nil {
				c.HTML(http.StatusBadRequest, "login-admin.html", gin.H{"err": err})
			}
			c.Redirect(302, "/login-admin")
		}
	})
	// 管理者ユーザーログイン
	router.POST("/login-admin", func(c *gin.Context) {
		// sessionを作成
		session := sessions.Default(c)
		session.Set("loginUser", c.PostForm("admin_id"))
		session.Save()
		log.Println(session.Get("loginUser"))

		// ログインしているAdminの取得
		admin := control.GetAdmin(c.PostForm("admin_id"))
		// DBから取得したユーザーパスワード(Hash)
		dbPassword := admin.Password
		// フォームから取得したユーザーパスワード
		formPassword := c.PostForm("password")

		// ユーザーパスワードの比較
		if err := crypto.CompareHashAndPassword(dbPassword, formPassword); err != nil {
			log.Println("Could not log in")
			c.HTML(http.StatusBadRequest, "login-admin.html", gin.H{"err": err})
			c.Abort()
		} else {
			log.Println("Could log in")
			c.Redirect(302, "/home-admin")
		}
	})
	// 管理者ユーザーホーム画面
	router.GET("/home-admin", func(c *gin.Context) {
		session := sessions.Default(c)
		// if !session {
		//     c.Redirect(302, "/login")
		// }
		session_id := session.Get("loginUser")
		admin_id := session_id.(string)
		// 各学科ごとに学生を全件取得
		students_network := control.GetAllStudent("network")
		students_information := control.GetAllStudent("information")
		students_system := control.GetAllStudent("system")

		// 各学科ごとに研究室を全件取得
		labs_network := control.GetAllLab("network")
		labs_information := control.GetAllLab("information")
		labs_system := control.GetAllLab("system")
		c.HTML(200, "home-admin.html", gin.H{
			"admin_id":             admin_id,
			"students_network":     students_network,
			"students_information": students_information,
			"students_system":      students_system,
			"labs_network":         labs_network,
			"labs_information":     labs_information,
			"labs_system":          labs_system,
		})
	})
	// 管理者ユーザー情報新規登録画面
	router.GET("/register", func(c *gin.Context) {

		c.HTML(200, "register.html", gin.H{})
	})
	// 管理者ユーザー情報削除画面
	router.GET("/delete", func(c *gin.Context) {

		c.HTML(200, "delete.html", gin.H{})
	})
	// 学生データの削除
	router.POST("delete-student", func(c *gin.Context) {
		student_id := c.PostForm("student_id")
		// 削除エラーの場合にログに表示
		if err := control.DeleteStudent(student_id); err != nil {
			c.Redirect(302, "/home-admin")
			log.Println(err)
		} else {
			c.Redirect(302, "/home-admin")
		}
	})
	// 研究室データの削除
	router.POST("delete-lab", func(c *gin.Context) {
		lab_id := c.PostForm("lab_id")
		// 削除エラーの場合にログに表示
		if err := control.DeleteLab(lab_id); err != nil {
			c.Redirect(302, "/home-admin")
			log.Println(err)
		} else {
			c.Redirect(302, "/home-admin")
		}
	})
	// 志望書データの削除
	router.POST("delete-aspire", func(c *gin.Context) {
		aspire_id_int, err := strconv.Atoi(c.PostForm("aspire_id"))
		if err != nil {
			log.Println(err)
		}
		log.Println(aspire_id_int)
		log.Printf("%T\n", aspire_id_int) // int
		// 削除エラーの場合にログに表示
		if err := control.DeleteAspire(aspire_id_int); err != nil {
			c.Redirect(302, "/home-admin")
			log.Println(err)
		} else {
			c.Redirect(302, "/home-admin")
		}
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
			c.HTML(http.StatusBadRequest, "home-admin.html", gin.H{"err": err})
			c.Abort()
		} else {
			student_id := c.PostForm("student_id")
			password := c.PostForm("password")
			department := c.PostForm("department")
			// 登録ユーザーが重複していた場合にはじく処理
			if err := control.CreateStudent(student_id, password, department); err != nil {
				c.Redirect(302, "/home-admin")
				log.Println(err)
			} else {
				c.Redirect(302, "/home-admin")
			}
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
	router.POST("/logout", func(c *gin.Context) {
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
			"labs":       labs,
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
			c.HTML(http.StatusBadRequest, "home-admin.html", gin.H{"err": err})
			c.Abort()
		} else {
			lab_id := c.PostForm("lab_id")
			password := c.PostForm("password")
			department := c.PostForm("department")
			// 登録ユーザーが重複していた場合にはじく処理(errがある場合とない場合で処理が分けられていない)
			if err := control.CreateLab(lab_id, password, department); err != nil {
				c.Redirect(302, "/home-admin")
				log.Println(err)
				// c.HTML(http.StatusBadRequest, "register.html", gin.H{"err": err})
			} else {
				c.Redirect(302, "/")
			}
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
			"lab_id":  lab_id,
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
