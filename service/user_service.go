package service

import (
	"errors"

	"github.com/onainadapdap1/dev/kode/my_gram/helpers"
	"github.com/onainadapdap1/dev/kode/my_gram/models"
	"github.com/onainadapdap1/dev/kode/my_gram/repository"
)

type UserServiceInterface interface {
	Registeruser(input models.RegisterUserInput) (models.User, error)
	LoginUser(input models.LoginUserInput) (models.User, error)
	GetUserByID(ID uint) (models.User, error)
}

type userService struct {
	repo repository.UserRepositoryInterface
}

func NewUserService(repo repository.UserRepositoryInterface) UserServiceInterface {
	return &userService{repo: repo}
}

func (s *userService) Registeruser(input models.RegisterUserInput) (models.User, error) {
	user := models.User {
		Username: input.Username,
		Email: input.Email,
		Password: input.Password,
		Age: input.Age,
	}

	newUser, err := s.repo.RegisterUser(user)
	if err != nil {
		return user, err
	}

	return newUser, nil
}

func (s *userService) LoginUser(input models.LoginUserInput) (models.User, error) {
	email := input.Email
	inputPassword := input.Password

	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return user, err
	}

	if user.ID == 0 {
		return user, errors.New("no user found on that email")
	}

	comparePass := helpers.ComparePassword([]byte(user.Password), []byte(inputPassword))
	if !comparePass {
		return user, err
	}

	return user, nil
}

func (s *userService) GetUserByID(ID uint) (models.User, error) {
	user, err := s.repo.FindByID(ID)
	if err != nil {
		return user, err
	}

	if user.ID == 0 {
		return user, errors.New("no user found on with that ID")
	}

	return user, nil
}