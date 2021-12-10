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
func CreateLab(lab_id string, password string, department string) []error {
	passwordEncrypt, _ := crypto.PasswordEncrypt(password)
	db := gormConnect()
	defer db.Close()
	// Insert処理
	if err := db.Create(&models.Lab{Lab_id: lab_id, Password: passwordEncrypt, Department: department}).GetErrors(); err != nil {
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

// 特定の学科に対応するLabを全件取得
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

// 志望研究室、理由、志望度をAspireに登録する処理
func CreateAspire(student_id string, lab_id string, reason string, rank string) {
	db := gormConnect()
	// Insert処理
	db.Create(&models.Aspire{Student_id: student_id, Lab_id: lab_id, Reason: reason, Rank: rank})
}
