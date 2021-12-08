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
// Aspireモデルの宣言
type Aspire struct {
	Student_id string `form:"student_id" binding:"required" gorm:"unique;not null"`
	Lab_id string `form:"lab_id" binding:"required" gorm:"unique;not null"`
	Reason string `form:"reason"`
	Rank uint8 `form:"lank"`
}