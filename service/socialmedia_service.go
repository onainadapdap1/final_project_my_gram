package service

import (
	"github.com/onainadapdap1/dev/kode/my_gram/dtos"
	"github.com/onainadapdap1/dev/kode/my_gram/models"
	"github.com/onainadapdap1/dev/kode/my_gram/repository"
)

type SocialmediaServiceInterface interface {
	CreateSocialMedia(input *dtos.CreateSocialMediaInput) (*models.SocialMedia, error)
	FindSocialMediaByID(inputID *dtos.GetSocialMediaDetailInput) (*models.SocialMedia, error)
	UpdateSocialMedia(inputID *dtos.GetSocialMediaDetailInput, inputData *dtos.CreateSocialMediaInput) (*models.SocialMedia, error)
	FindAllSocialMedia() (*[]models.SocialMedia, error)
	DeleteSocialMedia(ID uint) error
	RestoreSocialMedia(uint) (*models.SocialMedia, error)
	GetDeletedSocialMediaByID(ID uint) (*models.SocialMedia, error)
}

type socialmediaService struct {
	repo repository.SocialmediaRepositoryInterface
}

func NewSocialmediaService(repository repository.SocialmediaRepositoryInterface) SocialmediaServiceInterface {
	return &socialmediaService{repo: repository}
}

func (s *socialmediaService) CreateSocialMedia(input *dtos.CreateSocialMediaInput) (*models.SocialMedia, error) {
	socialMedia := &models.SocialMedia{
		Name:           input.Name,
		SocialMediaUrl: input.SocialMediaUrl,
		UserID:         input.User.ID,
	}

	socialMedia, err := s.repo.CreateSocialMedia(socialMedia)
	if err != nil {
		return nil, err
	}

	return socialMedia, nil
}

func (s *socialmediaService) FindSocialMediaByID(inputID *dtos.GetSocialMediaDetailInput) (*models.SocialMedia, error) {
	socialMedia, err := s.repo.FindSocialMediaByID(inputID.ID)
	if err != nil {
		return nil, err
	}

	return socialMedia, nil
}

func (s *socialmediaService) UpdateSocialMedia(inputID *dtos.GetSocialMediaDetailInput, inputData *dtos.CreateSocialMediaInput) (*models.SocialMedia, error) {
	socialMedia, err := s.repo.FindSocialMediaByID(inputID.ID)
	if err != nil {
		return nil, err
	}

	socialMedia.Name = inputData.Name
	socialMedia.SocialMediaUrl = inputData.SocialMediaUrl

	updatedSocialMedia, err := s.repo.UpdateSocialMedia(socialMedia)
	if err != nil {
		return nil, err
	}

	return updatedSocialMedia, err
}

func (s *socialmediaService) FindAllSocialMedia() (*[]models.SocialMedia, error) {
	socialMedias, err := s.repo.FindAllSocialMedia()
	if err != nil {
		return nil, err
	}

	return socialMedias, nil
}

func (s *socialmediaService) DeleteSocialMedia(ID uint) error {
	socialMedia, err := s.repo.FindSocialMediaByID(ID)
	if err != nil {
		return err
	}

	err = s.repo.DeleteSocialMedia(socialMedia)
	if err != nil {
		return err
	}

	return nil
}

func (s *socialmediaService) GetDeletedSocialMediaByID(ID uint) (*models.SocialMedia, error) {
	socialMedia, err := s.repo.GetDeletedSocialMediaByID(ID)
	if err != nil {
		return nil, err
	}

	return socialMedia, nil
}

func (s *socialmediaService) RestoreSocialMedia(ID uint) (*models.SocialMedia, error) {
	socialMedia, err := s.repo.FindSocialMediaByID(ID)
	if err != nil {
		return nil, err
	}

	restoreSocialMedia, err := s.repo.RestoreSocialMedia(socialMedia)
	if err != nil {
		return nil, err
	}

	return restoreSocialMedia, nil
}
