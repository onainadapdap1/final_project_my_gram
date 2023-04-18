package repository

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/onainadapdap1/dev/kode/my_gram/models"
)

type UserRepositoryInterface interface {
	RegisterUser(user models.User) (models.User, error)
	FindByEmail(email string) (models.User, error)
	FindByID(ID uint) (models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepositoryInterface {
	return &userRepository{db: db}
}

func (r *userRepository) RegisterUser(user models.User) (models.User, error) {
	tx := r.db.Begin()

	if err := tx.Debug().Create(&user).Error; err != nil {
		tx.Rollback()
		return user, fmt.Errorf("[RegisterAdminRepoImpl.Insert] Error when query save data with : %w", err)
	}
	tx.Commit()

	return user, nil
}


func (r *userRepository) FindByEmail(email string) (models.User, error) {
	tx := r.db.Begin()
	var user models.User

	err := tx.Debug().Where("email = ?", email).First(&user).Error
	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *userRepository) FindByID(ID uint) (models.User, error) {
	tx := r.db.Begin()
	var user models.User

	if err := tx.Debug().Where("id = ?", ID).First(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}