package repository

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	"github.com/onainadapdap1/dev/kode/my_gram/models"
)

type SocialmediaRepositoryInterface interface {
	CreateSocialMedia(socialMedia *models.SocialMedia) (*models.SocialMedia, error)
	FindSocialMediaByID(ID uint) (*models.SocialMedia, error)
	UpdateSocialMedia(socialMedia *models.SocialMedia) (*models.SocialMedia, error)
	FindAllSocialMedia() (*[]models.SocialMedia, error)
	DeleteSocialMedia(*models.SocialMedia) error
	RestoreSocialMedia(*models.SocialMedia) (*models.SocialMedia, error)
	GetDeletedSocialMediaByID(ID uint) (*models.SocialMedia, error)
}

type socialmediaRepository struct {
	db *gorm.DB
}

func NewSocialmediaRepository(db *gorm.DB) SocialmediaRepositoryInterface {
	return &socialmediaRepository{db: db}
}

func (r *socialmediaRepository) CreateSocialMedia(socialMedia *models.SocialMedia) (*models.SocialMedia, error) {
	tx := r.db.Begin()

	if err := tx.Debug().Create(socialMedia).Error; err != nil {
		log.Println("[CreateSocialmedia.Insert] Error when query save data with:", err)
		tx.Rollback()
		return nil, fmt.Errorf("[CreateSocialmedia.Insert] Error when query save data with : %w", err)
	}
	tx.Commit()
	return socialMedia, nil
}

func (r *socialmediaRepository) FindSocialMediaByID(ID uint) (*models.SocialMedia, error) {
	tx := r.db.Begin()
	var socialMedia models.SocialMedia

	if err := tx.Debug().Unscoped().Preload("User").Where("id = ?", ID).First(&socialMedia).Error; err != nil {
		return nil, err
	}
	tx.Commit()

	return &socialMedia, nil
}

func (r *socialmediaRepository) UpdateSocialMedia(socialMedia *models.SocialMedia) (*models.SocialMedia, error) {
	tx := r.db.Begin()

	if err := tx.Debug().Save(socialMedia).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return socialMedia, nil
}

func (r *socialmediaRepository) FindAllSocialMedia() (*[]models.SocialMedia, error) {
	tx := r.db.Begin()
	socialMedias := []models.SocialMedia{}
	if err := tx.Debug().Preload("User").Find(&socialMedias).Error; err != nil {
		return nil, err
	}
	tx.Commit()

	return &socialMedias, nil
}

func (r *socialmediaRepository) DeleteSocialMedia(socialMedia *models.SocialMedia) error {
	tx := r.db.Begin()

	if err := tx.Debug().Delete(socialMedia).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}

func (r *socialmediaRepository) GetDeletedSocialMediaByID(ID uint) (*models.SocialMedia, error) {
	tx := r.db.Begin()
	socialMedia := models.SocialMedia{}

	if err := tx.Debug().Unscoped().Preload("User").Where("id = ?", ID).First(&socialMedia).Error; err != nil {
		return &socialMedia, fmt.Errorf("[GetByID.Get] Error when query get data with : %w", err)
	}

	tx.Commit()

	return &socialMedia, nil

}

func (r *socialmediaRepository) RestoreSocialMedia(socialMedia *models.SocialMedia) (*models.SocialMedia, error) {
	tx := r.db.Begin()

	if err := tx.Debug().Model(&models.SocialMedia{}).Unscoped().UpdateColumn("deleted_at", nil).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()

	return socialMedia, nil
}
