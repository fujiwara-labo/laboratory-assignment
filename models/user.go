package models

import (
	"github.com/jinzhu/gorm"
)

// Student モデルの宣言
type Student struct {
	gorm.Model
	Student_id string `form:"student_id" binding:"required" gorm:"unique;not null"`
	Password   string `form:"password" binding:"required"`
	Department string `form:"department" binding:"required"`
	Assign_lab string
}

// Labモデルの宣言
type Lab struct {
	gorm.Model
	Lab_id      string `form:"lab_id" binding:"required" gorm:"unique;not null"`
	Password    string `form:"password" binding:"required"`
	Department  string `form:"department" binding:"required"`
	Assign_max  int
	Assign_flag bool
}

// Aspireモデルの宣言(なぜかaspire_idは自動インクリメントになっている)
type Aspire struct {
	Aspire_id  int    `gorm:"primary_key;AUTO_INCREMENT"`
	Student_id string `form:"student_id" binding:"required"`
	Lab_id     string `form:"lab_id" binding:"required"`
	Reason     string `form:"reason"`
}

// Admin モデルの宣言
type Admin struct {
	gorm.Model
	Admin_id string `form:"admin_id" binding:"required" gorm:"unique;not null"`
	Password string `form:"password" binding:"required"`
}
