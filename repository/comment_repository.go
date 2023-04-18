package repository

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/onainadapdap1/dev/kode/my_gram/models"
)

type CommentRepositoryInterface interface {
	CreateComment(comment models.Comment) (models.Comment, error)
	FindPhotoByID(ID uint) (models.Photo, error)
	FindCommentByID(ID uint) (models.Comment, error)
	GetCommentByID(ID uint) (models.Comment, error)
	UpdateComment(comment models.Comment) (models.Comment, error)
	FindAllComments() ([]models.Comment, error)
	DeleteCommentByID(comment *models.Comment) error
	RestoreCommentByID(comment models.Comment) (models.Comment, error)
}

type commentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) CommentRepositoryInterface {
	return &commentRepository{db: db}
}

func (r *commentRepository) CreateComment(comment models.Comment) (models.Comment, error) {
	tx := r.db.Begin()
	if err := tx.Debug().Create(&comment).Error; err != nil {
		tx.Rollback()
		return comment, fmt.Errorf("[CreateComment.Insert] Error when query save data with : %w", err)
	}
	tx.Commit()

	return comment, nil
}

func (r *commentRepository) FindPhotoByID(ID uint) (models.Photo, error) {
	var photo models.Photo
	tx := r.db.Begin()

	if err := tx.Debug().Preload("User").Where("id = ?", ID).First(&photo).Error; err != nil {
		return photo, fmt.Errorf("[FindByID.Get] Error when query get data with : %w", err)
	}
	tx.Commit()

	return photo, nil
}

func (r *commentRepository) FindCommentByID(ID uint) (models.Comment, error) {
	tx := r.db.Begin()

	var comment models.Comment
	// user := models.User{}

	if err := tx.Debug().Preload("User").Preload("Photo").Where("id = ?", ID).Find(&comment).Error; err != nil {
		return comment, fmt.Errorf("[FindByID.Get] Error when query get data with : %w", err)
	}

	tx.Commit()
	return comment, nil
}

func(r *commentRepository) GetCommentByID(ID uint) (models.Comment, error) {
	tx := r.db.Begin()
	comment := models.Comment{}

	if err := tx.Debug().Where("id = ?", ID).First(&comment).Error; err != nil {
		return comment, fmt.Errorf("[GetByID.Get] Error when query get data with : %w", err)
	}
	tx.Commit()
	return comment, nil
}

func (r *commentRepository) UpdateComment(comment models.Comment) (models.Comment, error) {
	tx := r.db.Begin()

	if err := tx.Debug().Save(&comment).Error; err != nil {
		tx.Rollback()
		return comment, err
	}

	tx.Commit()

	return comment, nil
}

func (r *commentRepository) FindAllComments() ([]models.Comment, error) {
	tx := r.db.Begin()
	comments := []models.Comment{}

	if err := tx.Debug().Preload("User").Preload("Photo").Find(&comments).Error; err != nil {
		return comments, fmt.Errorf("[FindAllComments.Get] Error when query get data with : %w", err)
	}
	tx.Commit()

	return comments, nil
}

func (r *commentRepository) DeleteCommentByID(comment *models.Comment) error {
	// tx := r.db.Begin()
	// comment := models.Comment{}
	
	// if err := tx.Debug().Model(&comment).Unscoped().Delete(&comment).Error; err != nil {
	// 	return err
	// }
	// return nil
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(comment).Error; err != nil {
			return err
		}

		return nil
	})
	
}

func (r *commentRepository) RestoreCommentByID(comment models.Comment) (models.Comment, error) {
	tx := r.db.Begin()

	if err := tx.Debug().Model(comment).UpdateColumn("deleted_at", nil).Error; err != nil {
		tx.Rollback()
		return comment, err
	}
	tx.Commit()
	return comment, nil
}