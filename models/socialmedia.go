package models

import (
	"github.com/asaskevich/govalidator"
	"github.com/jinzhu/gorm"
)

type SocialMedia struct {
	gorm.Model
	Name           string `gorm:"not null" json:"name" form:"name" valid:"required~Your Name is required"`
	SocialMediaUrl string `gorm:"not null" json:"social_media_url" form:"social_media_url" valid:"required~Your Social Media Url is required"`
	UserID         uint
	User           User `valid:"-"`
}

func (social *SocialMedia) TableName() string {
	return "tb_socialmedias"
}

func (social *SocialMedia) BeforeSave(tx *gorm.DB) (err error) {
	_, errCreate := govalidator.ValidateStruct(social)
	if errCreate != nil {
		err = errCreate
		return
	}

	return
}
