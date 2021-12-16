package control

import (
	"log"
	"os"

	"github.com/fujiwara-labo/laboratory-assignment.git/crypto"
	"github.com/fujiwara-labo/laboratory-assignment.git/models"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
)

func gormConnect() *gorm.DB {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	DBMS := os.Getenv("DRIVER")
	CONNECT := os.Getenv("DSN")
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
	//構造体に基づいてテーブルを作成
	db.AutoMigrate(&models.Student{})
	log.Println("create Student table")
	db.AutoMigrate(&models.Lab{})
	log.Println("create Lab table")
	db.AutoMigrate(&models.Aspire{})
	log.Println("create Aspire table")
	db.AutoMigrate(&models.Admin{})
	log.Println("create Admin table")
}

// 学生ユーザー登録処理
func CreateStudent(student_id string, password string, department string) []error {
	passwordEncrypt, _ := crypto.PasswordEncrypt(password)
	db := gormConnect()
	defer db.Close()
	// Insert処理
	if err := db.Create(&models.Student{Student_id: student_id, Password: passwordEncrypt, Department: department, Assign_lab: "none"}).GetErrors(); err != nil {
		return err
	}
	return nil
}

// 教員ユーザー登録処理
func CreateLab(lab_id string, password string, department string, assign_max int) []error {
	passwordEncrypt, _ := crypto.PasswordEncrypt(password)
	db := gormConnect()
	defer db.Close()
	// Insert処理
	if err := db.Create(&models.Lab{Lab_id: lab_id, Password: passwordEncrypt, Department: department, Assign_max: assign_max, Assign_flag: false}).GetErrors(); err != nil {
		return err
	}
	return nil
}

// 管理者ユーザー登録処理
func CreateAdmin(admin_id string, password string) []error {
	passwordEncrypt, _ := crypto.PasswordEncrypt(password)
	db := gormConnect()
	defer db.Close()
	// Insert処理
	if err := db.Create(&models.Admin{Admin_id: admin_id, Password: passwordEncrypt}).GetErrors(); err != nil {
		return err
	}
	return nil
}

// 学生ユーザーを一件取得
func GetStudent(student_id string) models.Student {
	db := gormConnect()
	var student models.Student
	db.First(&student, "student_id = ?", student_id)
	db.Close()
	return student
}

// 教員ユーザーを一件取得
func GetLab(lab_id string) models.Lab {
	db := gormConnect()
	var lab models.Lab
	db.First(&lab, "lab_id = ?", lab_id)
	db.Close()
	return lab
}

// 管理者ユーザーを一件取得
func GetAdmin(admin_id string) models.Admin {
	db := gormConnect()
	var admin models.Admin
	db.First(&admin, "admin_id = ?", admin_id)
	db.Close()
	return admin
}

// 特定の学科に対応するStudentを全件取得
func GetAllStudent(department string) []models.Student {
	db := gormConnect()
	var students []models.Student
	db.Where("department = ?", department).Find(&students)
	db.Close()
	return students
}

// 特定の学科に対応するLabを全件取得
func GetAllLab(department string) []models.Lab {
	db := gormConnect()
	var labs []models.Lab
	db.Where("department = ?", department).Find(&labs)
	db.Close()
	return labs
}

// 志望書一覧を取得
func GetAspires() []models.Aspire {
	db := gormConnect()
	var aspires []models.Aspire
	db.Find(&aspires)
	db.Close()
	return aspires
}

// ログインしている研究室の志望書一覧を取得
func GetAllAspire(lab_id string) []models.Aspire {
	db := gormConnect()
	var aspires []models.Aspire
	db.Where("lab_id = ?", lab_id).Find(&aspires)
	db.Close()
	return aspires
}

// 各学生の提出した志望書一覧を取得
func GetSubmitAsp(student_id string) []models.Aspire {
	db := gormConnect()
	var aspires []models.Aspire
	db.Where("student_id = ?", student_id).Find(&aspires)
	db.Close()
	return aspires
}

// 各学生の提出した志望書数を取得
func GetSubmitAspNum(student_id string) int {
	db := gormConnect()
	var aspires []models.Aspire
	db.Where("student_id = ?", student_id).Find(&aspires)
	db.Close()
	submit_num := len(aspires)
	return submit_num
}

// Studentからassign_labがlab_idに一致する学生を全件取得
func GetAllAssignStudent(lab_id string) []models.Student {
	db := gormConnect()
	var students []models.Student
	db.Where("assign_lab = ?", lab_id).Find(&students)
	return students
}

// 志望研究室、理由、志望度をAspireに登録する処理
func CreateAspire(student_id string, lab_id string, reason string) {
	db := gormConnect()
	// Insert処理
	db.Create(&models.Aspire{Student_id: student_id, Lab_id: lab_id, Reason: reason})
}

// 同じ研究室に志望書を出していないか確認
func ConfExistSameAsp(student_id string, lab_id string, aspires []models.Aspire) bool {
	for _, aspire := range aspires {
		if lab_id == aspire.Lab_id {
			return true
		}
	}
	return false
}

// student_idに対応する学生の削除
func DeleteStudent(student_id string) []error {
	db := gormConnect()
	var student models.Student
	// delete処理
	if err := db.Where("student_id = ?", student_id).Unscoped().Delete(&student).GetErrors(); err != nil {
		return err
	}
	return nil
}

// lab_idに対応する研究室の削除
func DeleteLab(lab_id string) []error {
	db := gormConnect()
	var lab models.Lab
	// delete処理
	if err := db.Where("lab_id = ?", lab_id).Unscoped().Delete(&lab).GetErrors(); err != nil {
		return err
	}
	return nil
}

// aspire_idに対応する研究室の削除
func DeleteAspire(aspire_id int) []error {
	db := gormConnect()
	var aspire models.Aspire
	// delete処理
	if err := db.Where("aspire_id = ?", aspire_id).Unscoped().Delete(&aspire).GetErrors(); err != nil {
		return err
	}
	return nil
}

// student_idに対応する任意のデータの変更
func FixStudent(student_id string, new_data string) []error {
	db := gormConnect()
	var student models.Student
	// fix
	if err := db.Model(&student).Where("student_id = ?", student_id).Update("department", new_data).GetErrors(); err != nil {
		return err
	}
	return nil
}

// lab_idに対応する任意のデータの変更
func FixLab(lab_id string, department string, assign_max int) []error {
	db := gormConnect()
	var lab models.Lab
	// fix
	if err := db.Model(&lab).Where("lab_id = ?", lab_id).Update("department", department).GetErrors(); err != nil {
		log.Println(err)
	}
	if err := db.Model(&lab).Where("lab_id = ?", lab_id).Update("assign_max", assign_max).GetErrors(); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// assign_flagがfalseのlabを全件取得
func GetAllFalseLab() []models.Lab {
	db := gormConnect()
	var labs []models.Lab
	db.Where("assign_flag = ?", false).Find(&labs)
	db.Close()
	return labs
}

// assign_flagがfalseのlabを学科ごとに取得
func GetAllFalseLabByDep(department string) []models.Lab {
	db := gormConnect()
	var labs []models.Lab
	db.Where(&models.Lab{Department: department, Assign_flag: false}).Find(&labs)
	db.Close()
	return labs
}

// Aspireからlab_idごとに提出学生数を取得する関数
func GetSubmitNum(lab_id string) int {
	db := gormConnect()
	var aspires []models.Aspire
	db.Where("lab_id = ?", lab_id).Find(&aspires)
	db.Close()
	submit_num := len(aspires)
	return submit_num
}

// assign_maxと（提出学生数＋配属決定学生数）を比較する関数
func CompMaxSubmit(lab_id string) bool {
	submit_num := len(GetAllAspire(lab_id))
	students_num := len(GetAllAssignStudent(lab_id))
	lab := GetLab(lab_id)
	if (submit_num + students_num) > lab.Assign_max {
		return true
	} else {
		return false
	}
}

// 任意の研究室に対してassign_maxと配属決定学生数を比較する関数
func CompMaxAssingStudent(lab_id string) bool {
	db := gormConnect()
	var students []models.Student
	lab := GetLab(lab_id)
	db.Where("assign_lab = ?", lab_id).Find(&students)
	log.Println("配属決定学生数")
	log.Println(len(students))
	log.Println("配属可能人数")
	log.Println(lab.Assign_max)

	if len(students) >= lab.Assign_max {
		// 配属決定学生が配属可能上限以上である
		return true
	}
	db.Close()
	return false
}

// aspireを論理削除する関数
func LogicDeleteAspire(student_id string) {
	db := gormConnect()
	var aspire []models.Aspire
	db.Where("student_id = ?", student_id).Delete(&aspire)

}

// 教員が受け入れたい学生を登録する（手動決定）
func AssignStudent(student_id string, lab_id string) {
	db := gormConnect()
	var student models.Student
	lab := GetLab(lab_id)
	log.Println(lab.Assign_flag)
	if !lab.Assign_flag {
		err := db.Model(&student).Where("student_id = ?", student_id).Update("assign_lab", lab_id).GetErrors()
		log.Println(err)
		LogicDeleteAspire(student_id)
	}
	SetAsssigFlag(lab_id)
	db.Close()
}

// 配属人数が定員人数の時assign_flagをtrueにする関数
func SetAsssigFlag(lab_id string) {
	db := gormConnect()
	var lab models.Lab
	get_lab := GetLab(lab_id)
	students := GetAllAssignStudent(lab_id)
	if get_lab.Assign_max == len(students) {
		db.Model(&lab).Where("lab_id = ?", lab_id).Update("assign_flag", true)
		log.Println(lab.Assign_flag)
		// 志望書が残っていたら削除する
		aspires := GetAllAspire(lab_id)
		for _, aspire := range aspires {
			LogicDeleteAspire(aspire.Student_id)
		}
	}
	db.Close()
}

// 定員割れしている研究室の自動配属決定
func AutoAssign(lab_id string) {
	db := gormConnect()
	var student models.Student
	aspires := GetAllAspire(lab_id)
	lab := GetLab(lab_id)
	if lab.Assign_max >= len(aspires) {
		for _, aspire := range aspires {
			db.Model(&student).Where("student_id = ?", aspire.Student_id).Update("assign_lab", lab_id)
			LogicDeleteAspire(aspire.Student_id)
		}
	}
	db.Close()
}

// to do : assign_labがない学生を定員割れしている研究室にランダム配属
// 学科ごとの配属研究室が決まっていない学生の取得
func GetNoAssginStudents() ([]models.Student, []models.Student, []models.Student) {
	db := gormConnect()
	var n_students []models.Student
	var i_students []models.Student
	var s_students []models.Student
	db.Where(&models.Student{Department: "network", Assign_lab: "none"}).Find(&n_students)
	db.Where(&models.Student{Department: "information", Assign_lab: "none"}).Find(&i_students)
	db.Where(&models.Student{Department: "system", Assign_lab: "none"}).Find(&s_students)
	db.Close()
	return n_students, i_students, s_students
}

// 学科ごとに定員割れしている研究室のID配列を取得
func NoLimitLabs() ([]string, []string, []string) {
	n_labs := GetAllFalseLabByDep("network")
	i_labs := GetAllFalseLabByDep("information")
	s_labs := GetAllFalseLabByDep("system")
	var n_lab_ids []string
	var i_lab_ids []string
	var s_lab_ids []string
	// networkのlab_idスライス
	for _, n_lab := range n_labs {
		assign_possible_num := n_lab.Assign_max - len(GetAllAssignStudent(n_lab.Lab_id))
		for i := 0; i < assign_possible_num; i++ {
			n_lab_ids = append(n_lab_ids, n_lab.Lab_id)
		}
	}
	// informationのlab_idスライス
	for _, i_lab := range i_labs {
		assign_possible_num := i_lab.Assign_max - len(GetAllAssignStudent(i_lab.Lab_id))
		for i := 0; i < assign_possible_num; i++ {
			i_lab_ids = append(i_lab_ids, i_lab.Lab_id)
		}
	}
	// systemのlab_idスライス
	for _, s_lab := range s_labs {
		assign_possible_num := s_lab.Assign_max - len(GetAllAssignStudent(s_lab.Lab_id))
		for i := 0; i < assign_possible_num; i++ {
			s_lab_ids = append(s_lab_ids, s_lab.Lab_id)
		}
	}
	return n_lab_ids, i_lab_ids, s_lab_ids
}

// 全体の実装：（assign_labがない学生を定員割れしている研究室にランダム配属）
func UndecidedAssignment() {
	db := gormConnect()
	// 定員割れしている研究室に志望する学生を全員配属決定にする
	labs := GetAllFalseLab()
	for _, lab := range labs {
		AutoAssign(lab.Lab_id)
	}
	// 学科ごとの配属研究室が決まっていない学生の取得
	n_students, i_students, s_students := GetNoAssginStudents()
	log.Println("志望書を出していない学生のリスト")
	log.Println(n_students)
	// 学科ごとに定員割れしている研究室のID配列を取得
	n_labs, i_labs, s_labs := NoLimitLabs()
	log.Println("配属数が定員割れしている研究室のリスト")
	log.Println(n_labs)
	// ランダム配属処理
	for idx, n_student := range n_students {
		db.Model(&n_student).Update("assign_lab", n_labs[idx])
		LogicDeleteAspire(n_student.Student_id)
	}

	for idx, i_student := range i_students {
		db.Model(&i_student).Update("assign_lab", i_labs[idx])
		LogicDeleteAspire(i_student.Student_id)
	}

	for idx, s_student := range s_students {
		db.Model(&s_student).Update("assign_lab", s_labs[idx])
		LogicDeleteAspire(s_student.Student_id)
	}
	db.Close()
}
