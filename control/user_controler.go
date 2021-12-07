package control

import(
	"log"
	"net/http"
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
	// USER := os.Getenv("MYSQL_USER")
	// PASS := os.Getenv("MYSQL_PASSWORD")
	// DBNAME := os.Getenv("MYSQL_DATABASE")
	// CONNECT := USER + ":" + PASS + "@/" + DBNAME + "?parseTime=true"
	db, err := gorm.Open(DBMS, CONNECT)
	if err != nil {
		panic(err.Error())
	}
	return db
}
// DBの初期化
func dbInit() {
	db := gormConnect()
	// コネクション解放
	defer db.Close()
    //構造体に基づいてテーブルを作成
	db.AutoMigrate(&models.User{})
}

// ユーザー登録処理
func createUser(username string, password string) []error {
	passwordEncrypt, _ := crypto.PasswordEncrypt(password)
	db := gormConnect()
	defer db.Close()
	// Insert処理
	if err := db.Create(&models.User{Username: username, Password: passwordEncrypt}).GetErrors(); err != nil {
		return err
	}
	return nil

}
// ユーザーを一件取得
func getUser(username string) models.User {
	db := gormConnect()
	var user models.User
	db.First(&user, "username = ?", username)
	db.Close()
	return user
}