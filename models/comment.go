package models

import (
	"github.com/asaskevich/govalidator"
	"github.com/jinzhu/gorm"
)

type Comment struct {
	gorm.Model
	UserID uint
	User User `valid:"-"`
	PhotoID uint 
	Photo Photo `valid:"-"`
	Message string `gorm:"not null" json:"title" form:"title" valid:"required~Your Title is required"`
}

func (c *Comment) TableName() string {
	return "tb_comments"
}

func (c *Comment) BeforeSave(tx *gorm.DB) (err error) {
	if _, errCreate := govalidator.ValidateStruct(c); errCreate != nil {
		err = errCreate
		return
	}

	return
}
