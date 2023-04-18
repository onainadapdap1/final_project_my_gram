package service

import (
	"github.com/onainadapdap1/dev/kode/my_gram/dtos"
	"github.com/onainadapdap1/dev/kode/my_gram/models"
	"github.com/onainadapdap1/dev/kode/my_gram/repository"
)

type CommentServiceInterface interface {
	CreateComment(input dtos.CreateCommentInput) (models.Comment, error)
	FindPhotoByID(ID uint) (models.Photo, error)
	FindCommentByID(ID uint) (models.Comment, error)
	GetCommentByID(ID uint) (models.Comment, error)
	UpdateComment(inputID dtos.GetCommentDetailInput, inputData dtos.UpdateCommentInput) (models.Comment, error)
	FindAllComments() ([]models.Comment, error)
	DeleteCommentByID(ID uint) error
	RestoreCommentByID(ID uint) (models.Comment, error)
	FindDeletedCommentByID(ID uint) (models.Comment, error)
}

type commentService struct {
	repo repository.CommentRepositoryInterface
}

func NewCommentService(repo repository.CommentRepositoryInterface) CommentServiceInterface {
	return &commentService{repo: repo}
}

func (s *commentService) CreateComment(input dtos.CreateCommentInput) (models.Comment, error) {
	comment := models.Comment {
		Message: input.Message,
		UserID: input.User.ID,
		PhotoID: input.PhotoID,
	}

	newComment, err := s.repo.CreateComment(comment)
	if err != nil {
		return newComment, err
	}

	return newComment, nil
}

func (s *commentService) FindPhotoByID(ID uint) (models.Photo, error) {
	photo, err := s.repo.FindPhotoByID(ID); 
	if err != nil {
		return photo, err
	}

	return photo, nil
}

func (s *commentService) FindCommentByID(ID uint) (models.Comment, error) {
	comment, err := s.repo.FindCommentByID(ID)
	if err != nil {
		return comment, err
	}

	return comment, nil
}

func (s *commentService) GetCommentByID(ID uint) (models.Comment, error) {
	comment, err := s.repo.GetCommentByID(ID)
	if err != nil {
		return comment, err
	}

	return comment, nil
}

func (s *commentService) UpdateComment(inputID dtos.GetCommentDetailInput, inputData dtos.UpdateCommentInput) (models.Comment, error) {
	comment, err := s.repo.FindCommentByID(inputID.ID)
	if err != nil {
		return comment, err
	}
	// comment := models.Comment{}

	comment.Message = inputData.Message
	// comment.PhotoID = inputData.PhotoID

	updatedComment, err := s.repo.UpdateComment(comment)
	if err != nil {
		return updatedComment, err
	}
	return updatedComment, nil

}

func (s *commentService) FindAllComments() ([]models.Comment, error) {
	comments, err := s.repo.FindAllComments()
	if err != nil {
		return comments, err
	}

	return comments, nil
}

func (s *commentService) DeleteCommentByID(ID uint) error {
	comment, err := s.repo.FindCommentByID(ID)
	if err != nil {
		return err
	}
	err = s.repo.DeleteCommentByID(&comment)
	if err != nil {
		return err
	}

	return nil
}

func (s *commentService) FindDeletedCommentByID(ID uint) (models.Comment, error) {
	deletedComment, err := s.repo.FindDeletedCommentByID(ID)
	if err != nil {
		return deletedComment, err
	}

	return deletedComment, nil
}

func (s *commentService) RestoreCommentByID(ID uint) (models.Comment, error) {
	comment, err := s.repo.FindDeletedCommentByID(ID)
	if err != nil {
		return comment, err
	}
	restoreComment, err := s.repo.RestoreCommentByID(comment)
	if err != nil {
		return comment, err
	}

	return restoreComment, nil
}
