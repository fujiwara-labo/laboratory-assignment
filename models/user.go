package models

import(
	"github.com/jinzhu/gorm"
)
// Student モデルの宣言
type Student struct {
	gorm.Model
	Student_id string `form:"student_id" binding:"required" gorm:"unique;not null"`
	Password string `form:"password" binding:"required"`
	Department string `form:"department" binding:"required"`
}
// Labモデルの宣言
type Lab struct {
	gorm.Model
	Lab_id string `form:"lab_id" binding:"required" gorm:"unique;not null"`
	Password string `form:"password" binding:"required"`
	Department string `form:"department" binding:"required"`
}
// Aspireモデルの宣言(なぜかaspire_idは自動インクリメントになっている)
type Aspire struct {
	Aspire_id int `gorm:"primary_key"`
	Student_id string `form:"student_id" binding:"required"`
	Lab_id string `form:"lab_id" binding:"required"`
	Reason string `form:"reason"`
	Rank string `form:"lank"`
}