package repository

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	"github.com/onainadapdap1/dev/kode/my_gram/models"
)

type PhotoRepositoryInterface interface {
	CreatePhoto(photo models.Photo) (models.Photo, error)
	FindByPhotoID(ID uint) (models.Photo, error)
	UpdateCategory(photo models.Photo) (models.Photo, error)
	FindAllPhoto() ([]models.Photo, error)
	DeletePhotoByID(photo models.Photo) error
}

type photoRepository struct {
	db *gorm.DB
}

func NewPhotoRepository(db *gorm.DB) PhotoRepositoryInterface {
	return &photoRepository{db: db}
}

func (r *photoRepository) CreatePhoto(photo models.Photo) (models.Photo, error) {
	tx := r.db.Begin()

	if err := tx.Debug().Create(&photo).Error; err != nil {
		log.Println("[CreatePhoto.Insert] Error when query save data with:", err)
		tx.Rollback()
		return photo, fmt.Errorf("[CreatePhoto.Insert] Error when query save data with : %w", err)
	}

	tx.Commit()
	return photo, nil
}

func (r *photoRepository) FindByPhotoID(ID uint) (models.Photo, error) {
	var photo models.Photo
	tx := r.db.Begin()
	if err := tx.Debug().Preload("User").Where("id = ?", ID).First(&photo).Error; err != nil {
		return photo, err
	}

	return photo, nil
}

func (r *photoRepository) UpdateCategory(photo models.Photo) (models.Photo, error) {
	tx := r.db.Begin()
	if err := tx.Debug().Save(&photo).Error; err != nil {
		tx.Rollback()
		return photo, err
	}

	return photo, nil
}

func (r *photoRepository) FindAllPhoto() ([]models.Photo, error) {
	var photos []models.Photo
	tx := r.db.Begin()
	if err := tx.Debug().Preload("User").Find(&photos).Error; err != nil {
		return photos, err
	}

	return photos, nil
}

func (r *photoRepository) DeletePhotoByID(photo models.Photo) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("photo_id = ?", photo.ID).Unscoped().Delete(&models.Comment{}).Error; err != nil {
			return err
		}

		if err := tx.Unscoped().Delete(photo).Error; err != nil {
			return err
		}

		return nil
	})
}