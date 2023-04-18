package dtos

import "github.com/onainadapdap1/dev/kode/my_gram/models"


type CreatePhotoInput struct {
	Title    string `gorm:"not null" json:"title" form:"title"`
	Caption  string `gorm:"null" json:"caption" form:"caption"`
	PhotoUrl string `gorm:"not null" json:"photo_url" form:"photo_url"`
	User     models.User
}

type GetPhotoDetailInput struct {
	ID uint `uri:"id" binding:"required"`
}

type PhotoFormatter struct {
	ID uint `json:"id"`
	UserID uint `json:"user_id"`
	Title string `json:"title"`
	Caption string `json:"caption"`
	PhotoUrl string `json:"photo_url"`
}

func FormatePhoto(photo models.Photo) PhotoFormatter {
	photoFormatter := PhotoFormatter {
		ID: photo.ID,
		UserID: photo.UserID,
		Title: photo.Title,
		Caption: photo.Caption,
		PhotoUrl: photo.PhotoUrl,
	}

	return photoFormatter
}

type PhotoDetailFormatter struct {
	ID uint `json:"id"`
	UserID uint `json:"user_id"`
	Title string `json:"title"`
	Caption string `json:"caption"`
	PhotoUrl string `json:"photo_url"`
	User PhotoUserFormatter
}

type PhotoUserFormatter struct {
	Username string `json:"username"`
	Age int `json:"age"`
}

func FormatePhotoDetail(photo models.Photo) PhotoDetailFormatter {
	photoDetailFormatter := PhotoDetailFormatter {
		ID: photo.ID,
		UserID: photo.ID,
		Title: photo.Title,
		Caption: photo.Caption,
		PhotoUrl: photo.PhotoUrl,
	}
	user := photo.User
	photoUserFormatter := PhotoUserFormatter {
		Username: user.Username,
		Age: user.Age,
	}

	photoDetailFormatter.User = photoUserFormatter

	return photoDetailFormatter
}

func FormatePhotoDetails(photos []models.Photo) []PhotoDetailFormatter {
	photosDetailFormatter := []PhotoDetailFormatter{}
	for _, photo := range photos {
		photoDetailFormatter := FormatePhotoDetail(photo)
		photosDetailFormatter = append(photosDetailFormatter, photoDetailFormatter)
	}

	return photosDetailFormatter
}