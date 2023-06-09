package dtos

import "github.com/onainadapdap1/dev/kode/my_gram/models"

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

type UserRegisterFormatter struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Age      int    `json:"age"`
}

// format user response register
func FormateUserRegister(user models.User) UserRegisterFormatter {
	formatter := UserRegisterFormatter{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Age:      user.Age,
	}

	return formatter
}

type LoginUserFormatter struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Age      int    `json:"age"`
	Token    string `json:"token"`
}

func FormateUserLogin(user models.User, token string) LoginUserFormatter {
	formatter := LoginUserFormatter{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Age:      user.Age,
		Token:    token,
	}

	return formatter
}
