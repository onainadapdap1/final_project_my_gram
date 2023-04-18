package models

import (
	"errors"
	"log"

	"github.com/asaskevich/govalidator"
	"github.com/jinzhu/gorm"
	"github.com/onainadapdap1/dev/kode/my_gram/helpers"
)

type User struct {
	gorm.Model
	Username string `gorm:"not null;uniqueIndex" json:"username" form:"username" valid:"required~Your Username is required"`
	Email    string `gorm:"not null;uniqueIndex" json:"email" form:"email" valid:"required~Your Email is required,email~Invalid Email format"`
	Password string `gorm:"not null" json:"password" form:"password" valid:"required~Your Password is required,minstringlength(6)~Your Password must be at least 6 characters"`
	Age      int    `gorm:"not null" json:"age" form:"age" valid:"required~Your Age is required"`
	Photo   []Photo
	Comment []Comment
}

type RegisterUserInput struct {
	Username string `gorm:"not null;uniqueIndex" json:"username" form:"username" `
	Email    string `gorm:"not null;uniqueIndex" json:"email" form:"email" `
	Password string `gorm:"not null" json:"password" form:"password" `
	Age      int    `gorm:"not null" json:"age" form:"age" `
}

type LoginUserInput struct {
	Email    string `gorm:"not null;uniqueIndex" json:"email" form:"email" `
	Password string `gorm:"not null" json:"password" form:"password" `
}

// naming convention
func (u *User) TableName() string {
	return "tb_users"
}

func (u *User) BeforeSave(tx *gorm.DB) (err error) {
	_, errCreate := govalidator.ValidateStruct(u)
	if u.Age < 8 {
		err = errors.New("you must be at least 8 years")
		return
	}
	if errCreate != nil {
		err = errCreate
		return
	}

	u.Password, err = helpers.HassPass(u.Password)
	if err != nil {
		log.Println("error while hashing password")
		return
	}
	err = nil
	return
}
