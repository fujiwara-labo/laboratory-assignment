package control

import(
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
	db.AutoMigrate(&models.User{})
}

// ユーザー登録処理
func CreateUser(username string, password string, department string) []error {
	passwordEncrypt, _ := crypto.PasswordEncrypt(password)
	db := gormConnect()
	defer db.Close()
	// Insert処理
	if err := db.Create(&models.User{Username: username, Password: passwordEncrypt, Department: department}).GetErrors(); err != nil {
		return err
	}
	return nil

}
// ユーザーを一件取得
func GetUser(username string) models.User {
	db := gormConnect()
	var user models.User
	db.First(&user, "username = ?", username)
	db.Close()
	return user
}