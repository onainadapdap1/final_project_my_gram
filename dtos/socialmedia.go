package dtos

import "github.com/onainadapdap1/dev/kode/my_gram/models"

type GetSocialMediaDetailInput struct {
	ID uint `uri:"id" binding:"required"`
}

type CreateSocialMediaInput struct {
	Name    string `gorm:"not null" json:"name" form:"name"`
	SocialMediaUrl  string `gorm:"null" json:"social_media_url" form:"social_media_url"`
	User models.User `gorm:"-" swaggerignore:"true"`
}

type SocialMediaFormatter struct {
	ID uint `json:"id"`
	UserID uint `json:"user_id"`
	Name string `json:"social_media_name"`
	SocialMediaUrl string `json:"social_media_url"`
}

func FormateSocialMedia(socialMedia *models.SocialMedia) *SocialMediaFormatter {
	socialMediaFormatter := &SocialMediaFormatter {
		ID: socialMedia.ID,
		UserID: socialMedia.UserID,
		Name: socialMedia.Name,
		SocialMediaUrl: socialMedia.SocialMediaUrl,
	}

	return socialMediaFormatter
}

type SocialMediaDetailFormatter struct {
	ID uint `json:"id"`
	UserID uint `json:"user_id"`
	Name string `json:"social_media_name"`
	SocialMediaUrl string `json:"social_media_url"`
	User SocialMediaUserFormatter
}

type SocialMediaUserFormatter struct {
	Username string `json:"username"`
	Age int `json:"age"`
}

func FormateSocialMediaDetail(socialMedia *models.SocialMedia) *SocialMediaDetailFormatter {
	socialMediaDetailFormatter := &SocialMediaDetailFormatter {
		ID: socialMedia.ID,
		UserID: socialMedia.UserID,
		Name: socialMedia.Name,
		SocialMediaUrl: socialMedia.SocialMediaUrl,
	}
	user := socialMedia.User
	socialMediaUserFormatter := SocialMediaUserFormatter {
		Username: user.Username,
		Age: user.Age,
	}

	socialMediaDetailFormatter.User = socialMediaUserFormatter

	return socialMediaDetailFormatter
}

func FormateSocialMediaDetails(socialMedias *[]models.SocialMedia) *[]SocialMediaDetailFormatter {
	socialMediasDetailFormatter := []SocialMediaDetailFormatter{}
	for _, socialMedia := range *socialMedias {
		socialMediaDetailFormatter := FormateSocialMediaDetail(&socialMedia)
		socialMediasDetailFormatter = append(socialMediasDetailFormatter, *socialMediaDetailFormatter)
	}

	return &socialMediasDetailFormatter
}