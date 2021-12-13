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
	if err := db.Create(&models.Student{Student_id: student_id, Password: passwordEncrypt, Department: department}).GetErrors(); err != nil {
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

// ログインしている研究室の志望書一覧を取得
func GetAllAspire(lab_id string) []models.Aspire {
	db := gormConnect()
	var aspires []models.Aspire
	db.Where("lab_id = ?", lab_id).Find(&aspires)
	db.Close()
	return aspires
}

// Studentからassign_labがlab_idに一致する学生を全件取得
func GetAllAssignStudent(lab_id string) []models.Student {
	db := gormConnect()
	var students []models.Student
	db.Where("assign_lab = ?", lab_id).Find(&students)
	return students
}

// 志望研究室、理由、志望度をAspireに登録する処理
func CreateAspire(student_id string, lab_id string, reason string, rank string) {
	db := gormConnect()
	// Insert処理
	db.Create(&models.Aspire{Student_id: student_id, Lab_id: lab_id, Reason: reason, Rank: rank})
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

// Aspireからlab_idごとに提出学生数を取得する関数
func GetSubmitNum(lab_id string) int {
	db := gormConnect()
	var aspires []models.Aspire
	db.Where("lab_id = ?", lab_id).Find(&aspires)
	db.Close()
	submit_num := len(aspires)
	return submit_num
}

// assign_maxと提出学生数を比較する関数
func CompMaxSubmit(submit_num int, assign_max int) bool {
	if submit_num > assign_max {
		return true
	} else {
		return false
	}
}

// 任意の研究室に対してassign_maxと配属決定学生数を比較する関数
func CompMaxAssingStudent(lab_id string) bool {
	db := gormConnect()
	var students []models.Student
	var lab models.Lab
	db.Where("assign_lab = ?", lab_id).Find(&students)

	if len(students) >= lab.Assign_max {
		// 配属決定学生が配属可能上限以上である
		return true
	} else {
		return false
	}
}

// 配属希望調査全体の関数
// to do : Studentのassign_labを決定する、希望数が多い場合は希望書を返す
func AssignResarch() {
	db := gormConnect()
	var students []models.Student
	var aspires []models.Aspire
	labs := GetAllFalseLab()

	for _, lab := range labs {
		submit_num := GetSubmitNum(lab.Lab_id)
		flag := CompMaxSubmit(submit_num, lab.Assign_max)

		if flag {
			// 希望学生数が配属上限数より多い場合 → false
			// 手動で配属学生を決定(AssignStudent)
		} else {
			// 希望学生数が配属上限数より少ない場合→配属決定(true)
			// 希望書を上記の条件を満たす研究室にいくつか提出している学生は複数の研究室に配属してしまう問題
			err := db.Where("lab_id = ?", lab.Lab_id).Find(&aspires).GetErrors()
			log.Println(err)
			for _, aspire := range aspires {
				log.Println(aspires)
				db.Model(&students).Where("student_id = ?", aspire.Student_id).Update("assign_lab", lab.Lab_id)
				db.Model(&lab).Where("lab_id = ?", lab.Lab_id).Update("assign_flag", true)
			}
		}
	}
	db.Close()
}

// 希望数が定員数を超えていた場合の配属決定処理の関数
func AssignStudent(student_id string, lab_id string) {
	db := gormConnect()
	var student models.Student
	var lab models.Lab

	// db.Model(&student).Where("student_id = ?", student_id).Update("assign_lab", lab_id)
	err := db.Model(&student).Where("student_id = ?", student_id).Update("assign_lab", lab_id).GetErrors()
	// db.Create(&models.Student{Assign_lab: lab_id})
	log.Println(err)
	flag := CompMaxAssingStudent(lab.Lab_id)
	if flag {
		db.Model(&lab).Where("lab_id = ?", lab_id).Update("assign_flag", true)
	}
	db.Close()
}
