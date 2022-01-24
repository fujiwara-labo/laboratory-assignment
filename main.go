package main

import (
	"github.com/fujiwara-labo/laboratory-assignment.git/control"
	"github.com/fujiwara-labo/laboratory-assignment.git/server"
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

	router.GET("/", server.GetHome())
	// 管理者ユーザー登録、ログイン画面
	router.GET("/login-admin", server.AdminLoginPage())
	// 管理者ユーザーログアウト
	router.POST("/logout-admin", server.AdminLogout())
	// 管理者ユーザー登録
	router.POST("/signup-admin", server.AdminRegister())
	// 管理者ユーザーログイン
	router.POST("/login-admin", server.AdminLogin())
	// 管理者ユーザーホーム画面
	router.GET("/home-admin", server.AdminHomePage())
	// ユーザー情報新規登録画面
	router.GET("/register", server.AdminUserRegisterPage())
	// ユーザー情報削除画面
	router.GET("/delete", server.AdminUserDeletePage())
	// 学生データの削除
	router.POST("delete-student", server.AdminStudentDelete())
	// 研究室データの削除
	router.POST("delete-lab", server.AdminLabDelete())
	// 志望書データの削除
	router.POST("delete-aspire", server.AdminAspireDelete())
	// 管理者ユーザー情報修正画面
	router.GET("/fix", server.AdminUserFixPage())
	// 学生データの変更
	router.POST("fix-student", server.AdminStudentFix())
	// Labデータの変更
	router.POST("fix-lab", server.AdminLabFix())
	// 学生ユーザー、ログイン画面
	router.GET("/login", server.StudentloginPage())
	// 学生ユーザー登録
	router.POST("/signup", server.StudentRegister())

	// 学生ユーザーログイン
	router.POST("/login", server.Studentlogin())
	// 学生ユーザーログアウト（sessionクリア）
	router.POST("/logout", server.Studentlogout())
	// 学生ユーザーホーム画面
	router.GET("/home-student", server.StudentHomePage())
	// 学生ユーザーの志望書提出フォーム画面
	router.GET("/form", server.AspireAdmitFormpage())
	// フォームの取得
	router.POST("/form", server.StudentAspireAdmit())

	// 教員ユーザー登録、ログイン画面
	router.GET("/login-lab", server.LabLogin())
	// 教員ユーザー登録
	router.POST("/signup-lab", server.LabRegister())
	// 研究室ページホーム
	router.GET("/home-lab", server.LabHomePage())
	// 配属学生決定ページ
	router.GET("/assign-lab", server.LabSelectStudentPage())
	// 教員ユーザーログイン
	router.POST("/login-lab", server.LabLoginPage())
	// 研究室配属先未決定者をランダム割り振り(研究室配属の自動決定) home-admin
	router.POST("/assign", server.AssignLab())
	// 研究室配属可能上限人数の自動設定　set-assign-max
	router.POST("/set-asssign-num", server.SetAssignMaxNum())
	// 研究室配属の手動決定 assign-lab
	router.POST("/select-students", server.AutoAssignLab())

	// 学生ユーザーの配属希望調査画面
	router.GET("/assign-reserch", server.AssignReserchPage())
	// 学生ユーザーの配属希望調査機能
	router.POST("/assign-reserch", server.AssignReserch())

	router.Run(":8080")

}
