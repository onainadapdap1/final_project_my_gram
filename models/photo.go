package models

import (
	"github.com/asaskevich/govalidator"
	"github.com/jinzhu/gorm"
)

type Photo struct {
	gorm.Model
	Title    string `gorm:"not null" json:"title" form:"title" valid:"required~Your Title is required"`
	Caption  string `gorm:"null" json:"caption" form:"caption"`
	PhotoUrl string `gorm:"not null" json:"photo_url" form:"photo_url" valid:"required~Your Photo Url is required"`
	UserID   uint
	User     User `valid:"-"`
	Comments []Comment `valid:"-"`
}


func (p *Photo) TableName() string {
	return "tb_photos"
}

func (p *Photo) BeforeSave(tx *gorm.DB) (err error) {
	_, errCreate := govalidator.ValidateStruct(p)

	if errCreate != nil {
		err = errCreate
		return
	}

	return
}
