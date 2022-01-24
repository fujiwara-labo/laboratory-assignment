package server

import (
	"log"
	"net/http"
	"strconv"

	// "github.com/fujiwara-labo/laboratory-assignment.git/control"

	"github.com/fujiwara-labo/laboratory-assignment.git/control"
	"github.com/fujiwara-labo/laboratory-assignment.git/crypto"
	"github.com/fujiwara-labo/laboratory-assignment.git/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// 学生教員選択画面
func GetHome() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			// htmlに渡す変数を定義
			"message": "test",
		})
	}
}

// // 管理者ユーザー登録、ログイン画面
func AdminLoginPage() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(200, "login-admin.html", gin.H{})
	}
}

// 管理者ユーザーログアウト機能
func AdminLogout() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		session.Clear()
		session.Save()
		log.Println(session.Get("loginUser"))
		c.Redirect(302, "/login-admin")
	}
}

// 管理者ユーザー登録機能
func AdminRegister() gin.HandlerFunc {
	return func(c *gin.Context) {
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
	}
}

// 管理者ユーザーログイン機能
func AdminLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
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
	}
}

// 管理者ユーザーホーム画面
func AdminHomePage() gin.HandlerFunc {
	return func(c *gin.Context) {
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

		// 志望書を全て取得
		aspires := control.GetAspires()
		c.HTML(200, "home-admin.html", gin.H{
			"admin_id":             admin_id,
			"students_network":     students_network,
			"students_information": students_information,
			"students_system":      students_system,
			"labs_network":         labs_network,
			"labs_information":     labs_information,
			"labs_system":          labs_system,
			"aspires":              aspires,
		})
	}
}

// ユーザー情報新規登録画面
func AdminUserRegisterPage() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(200, "register.html", gin.H{})
	}
}

// ユーザー情報削除画面
func AdminUserDeletePage() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(200, "delete.html", gin.H{})
	}
}

// 学生データの削除機能
func AdminStudentDelete() gin.HandlerFunc {
	return func(c *gin.Context) {
		student_id := c.PostForm("student_id")
		// 削除エラーの場合にログに表示
		if err := control.DeleteStudent(student_id); err != nil {
			c.Redirect(302, "/home-admin")
			log.Println(err)
		} else {
			c.Redirect(302, "/home-admin")
		}
	}
}

// 研究室データの削除
func AdminLabDelete() gin.HandlerFunc {
	return func(c *gin.Context) {
		lab_id := c.PostForm("lab_id")
		// 削除エラーの場合にログに表示
		if err := control.DeleteLab(lab_id); err != nil {
			c.Redirect(302, "/home-admin")
			log.Println(err)
		} else {
			c.Redirect(302, "/home-admin")
		}
	}
}

// 志望書データの削除
func AdminAspireDelete() gin.HandlerFunc {
	return func(c *gin.Context) {
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
	}
}

// 管理者ユーザー情報修正画面
func AdminUserFixPage() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(200, "fix.html", gin.H{})
	}
}

// 学生データの変更
func AdminStudentFix() gin.HandlerFunc {
	return func(c *gin.Context) {
		student_id := c.PostForm("student_id")
		new_department := c.PostForm("department")
		// 修正エラーの場合にログに表示
		if err := control.FixStudent(student_id, new_department); err != nil {
			log.Println(err)
			c.Redirect(302, "/home-admin")
		}
		c.Redirect(302, "/home-admin")
	}
}

// Labデータの変更
func AdminLabFix() gin.HandlerFunc {
	return func(c *gin.Context) {
		lab_id := c.PostForm("lab_id")
		new_department := c.PostForm("department")
		assign_max_int, err := strconv.Atoi(c.PostForm("assign_max"))
		if err != nil {
			log.Println(err)
		}
		// 修正エラーの場合にログに表示
		if err := control.FixLab(lab_id, new_department, assign_max_int); err != nil {
			log.Println(err)
			c.Redirect(302, "/home-admin")
		}
		c.Redirect(302, "/home-admin")
	}
}

// 学生ユーザ、ログイン画面
func StudentloginPage() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(200, "login.html", gin.H{})
	}
}

// 学生ユーザー登録
func StudentRegister() gin.HandlerFunc {
	return func(c *gin.Context) {
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
	}
}

// 学生ユーザーログイン
func Studentlogin() gin.HandlerFunc {
	return func(c *gin.Context) {
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
	}
}

// 学生ユーザーログアウト（sessionクリア）
func Studentlogout() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		session.Clear()
		session.Save()
		log.Println(session.Get("loginUser"))
		c.Redirect(302, "/")
	}
}

// 学生ユーザーホーム画面
func StudentHomePage() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		// if !session {
		//     c.Redirect(302, "/login")
		// }
		session_id := session.Get("loginUser")
		student_id := session_id.(string)
		submit_num := control.GetSubmitAspNum(student_id)
		aspires := control.GetSubmitAsp(student_id)
		get_student := control.GetStudent(student_id)
		assign_lab := get_student.Assign_lab
		text := "配属未決定です。志望書を提出してください"
		if submit_num == 1 {
			text = "申請中です"
		}
		if len(assign_lab) == 1 {
			text = "配属が決定しました"
		}
		c.HTML(200, "home-student.html", gin.H{
			"student_id": student_id,
			"lab_id":     assign_lab,
			"submit_num": submit_num,
			"message":    text,
			"aspires":    aspires,
		})
	}
}

// 学生ユーザーの志望書提出フォーム画面
func AspireAdmitFormpage() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)

		session_id := session.Get("loginUser")
		student_id := session_id.(string)
		student := control.GetStudent(student_id)
		labs := control.GetAllLab(student.Department)

		submit_num := control.GetSubmitAspNum(student_id)
		message := "志望書を提出してください"
		if submit_num == 1 {
			message = "提出上限：これ以上提出できません"
		}
		c.HTML(200, "form.html", gin.H{
			"student_id": student.Student_id,
			"submit_num": 1 - submit_num,
			"message":    message,
			"department": student.Department,
			"labs":       labs,
		})
	}
}

// フォームの提出
func StudentAspireAdmit() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)

		session_id := session.Get("loginUser")
		student_id := session_id.(string)
		log.Println(student_id)
		submit_num := control.GetSubmitAspNum(student_id)
		if submit_num == 0 {
			reason := c.PostForm("reason")
			lab_id := c.PostForm("lab_id")
			log.Println(lab_id)
			control.CreateAspire(student_id, lab_id, reason)
			c.Redirect(302, "/home-student")
		} else {
			c.Redirect(302, "/home-student")
		}
	}
}

// 教員ユーザーログイン画面
func LabLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(200, "login-lab.html", gin.H{})
	}
}

// 教員ユーザー登録
func LabRegister() gin.HandlerFunc {
	return func(c *gin.Context) {
		var form models.Lab
		// バリデーション処理
		if err := c.Bind(&form); err != nil {
			c.HTML(http.StatusBadRequest, "home-admin.html", gin.H{"err": err})
			c.Abort()
		} else {
			lab_id := c.PostForm("lab_id")
			password := c.PostForm("password")
			department := c.PostForm("department")
			assign_max_int, err := strconv.Atoi(c.PostForm("assign_max"))
			if err != nil {
				log.Println(err)
			}
			// 登録ユーザーが重複していた場合にはじく処理(errがある場合とない場合で処理が分けられていない)
			if err := control.CreateLab(lab_id, password, department, assign_max_int); err != nil {
				c.Redirect(302, "/home-admin")
				log.Println(err)
				// c.HTML(http.StatusBadRequest, "register.html", gin.H{"err": err})
			} else {
				c.Redirect(302, "/")
			}
		}
	}
}

// 研究室ページホーム
func LabHomePage() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		// if !session {
		//     c.Redirect(302, "/login-lab")
		// }
		session_id := session.Get("loginUser")
		lab_id := session_id.(string)
		// Studentsからlab_idで配属が決定した学生を取得
		students := control.GetAllAssignStudent(lab_id)
		// 志望書提出数が定員を超えたことを確認
		flag_asp := control.CompMaxSubmit(lab_id)
		// 配属決定学生数が定員まで決定したことを確認
		flag_assign := control.CompMaxAssingStudent(lab_id)
		text := "配属希望学生が定員まで達していません"
		if flag_asp {
			text = "配属希望学生が定員数を超えました。採用学生を配属学生選択から決定してください"

		} else {
			text = "配属希望学生が定員まで達していません"
		}
		if flag_assign {
			text = "配属学生が決定しました"
		}
		c.HTML(200, "home-lab.html", gin.H{
			"lab_id":   lab_id,
			"message":  text,
			"students": students,
		})
	}
}

// 配属学生決定ページ
func LabSelectStudentPage() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		session_id := session.Get("loginUser")
		lab_id := session_id.(string)
		lab := control.GetLab(lab_id)
		assignd_num := len(control.GetAllAssignStudent(lab_id))
		// 提出された志望書一覧を取得
		aspires := control.GetAllAspire(lab_id)
		c.HTML(200, "assign-lab.html", gin.H{
			"lab_id":    lab_id,
			"assin_num": lab.Assign_max - assignd_num,
			"lab_id2":   lab_id,
			"aspires":   aspires,
		})
	}
}

// 教員ユーザーログイン
func LabLoginPage() gin.HandlerFunc {
	return func(c *gin.Context) {
		// sessionを作成
		session := sessions.Default(c)
		session.Set("loginUser", c.PostForm("lab_id"))
		session.Save()
		// ログインしているLabの取得
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
	}
}

// 研究室配属先未決定者をランダム割り振り(研究室配属の自動決定) home-admin
func AssignLab() gin.HandlerFunc {
	return func(c *gin.Context) {
		control.UndecidedAssignment()
		c.Redirect(302, "/home-admin")
	}
}

// 研究室配属の手動決定 assign-lab
func AutoAssignLab() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		session_id := session.Get("loginUser")
		lab_id := session_id.(string)

		student_id := c.PostForm("student_id")
		control.AssignStudent(student_id, lab_id)
		c.Redirect(302, "/home-lab")
	}
}

func SetAssignMaxNum() gin.HandlerFunc {
	return func(c *gin.Context) {
		control.SetAssignMax()
		c.Redirect(302, "/home-admin")
	}
}

// 配属希望調査の画面表示
func AssignReserchPage() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)

		session_id := session.Get("loginUser")
		student_id := session_id.(string)
		student := control.GetStudent(student_id)
		labs := control.GetAllLab(student.Department)

		c.HTML(200, "assign-reserch.html", gin.H{
			"student_id": student_id,
			"department": student.Department,
			"labs":       labs,
		})
	}
}

func AssignReserch() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		session_id := session.Get("loginUser")
		student_id := session_id.(string)
		lab1 := c.PostForm("lab_id_1")
		lab2 := c.PostForm("lab_id_2")
		lab3 := c.PostForm("lab_id_3")
		log.Println(student_id)
		log.Println(lab1)
		log.Println(lab2)
		log.Println(lab3)

		c.Redirect(302, "/home-student")
	}
}
