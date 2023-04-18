package service

import (
	"github.com/onainadapdap1/dev/kode/my_gram/dtos"
	"github.com/onainadapdap1/dev/kode/my_gram/models"
	"github.com/onainadapdap1/dev/kode/my_gram/repository"
)

type PhotoServiceInterface interface {
	CreatePhoto(input dtos.CreatePhotoInput) (models.Photo, error)
	UpdatePhoto(inputID dtos.GetPhotoDetailInput, inputData dtos.CreatePhotoInput) (models.Photo, error)
	FindByPhotoID(ID uint) (models.Photo, error)
	FindAllPhoto() ([]models.Photo, error)
	DeletePhotoByID(photo models.Photo) error
	
}

type photoService struct {
	repo repository.PhotoRepositoryInterface
}

func NewPhotoService(repository repository.PhotoRepositoryInterface) PhotoServiceInterface {
	return &photoService{repo: repository}
}

func (s *photoService) CreatePhoto(input dtos.CreatePhotoInput) (models.Photo, error) {
	photo := models.Photo{
		UserID:   input.User.ID,
		Title:    input.Title,
		Caption:  input.Caption,
		PhotoUrl: input.PhotoUrl,
	}

	photo, err := s.repo.CreatePhoto(photo)
	if err != nil {
		return photo, err
	}

	return photo, nil
}

func (s *photoService) UpdatePhoto(inputID dtos.GetPhotoDetailInput, inputData dtos.CreatePhotoInput) (models.Photo, error) {
	photo, err := s.repo.FindByPhotoID(inputID.ID)
	if err != nil {
		return photo, err
	}

	photo.Title = inputData.Title
	photo.Caption = inputData.Caption
	photo.PhotoUrl = inputData.PhotoUrl

	updatedPhoto, err := s.repo.UpdateCategory(photo)
	if err != nil {
		return updatedPhoto, err
	}

	return updatedPhoto, nil
}

func (s *photoService) FindByPhotoID(ID uint) (models.Photo, error) {
	photo, err := s.repo.FindByPhotoID(ID)
	if err != nil {
		return photo, err
	}

	return photo, nil

}

func (s *photoService) FindAllPhoto() ([]models.Photo, error) {
	photos, err := s.repo.FindAllPhoto()
	if err != nil {
		return photos, err
	}

	return photos, nil
}

func (s *photoService) DeletePhotoByID(photo models.Photo) error {
	if err := s.repo.DeletePhotoByID(photo); err != nil {
		return err
	}

	return nil
}