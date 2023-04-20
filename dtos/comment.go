package dtos

import (
	"github.com/onainadapdap1/dev/kode/my_gram/models"
)

type GetCommentDetailInput struct {
	ID uint `uri:"id" binding:"required"`
}

type CreateCommentInput struct {
	Message string       `gorm:"not null" form:"message" json:"message" binding:"required"`
	PhotoID uint         `gorm:"not null" form:"photo_id" json:"photo_id" binding:"required"`
	User    models.User  `gorm:"-" swaggerignore:"true"`
	Photo   models.Photo `gorm:"-" swaggerignore:"true"`
}

type UpdateCommentInput struct {
	Message string       `gorm:"not null" form:"message" json:"message" binding:"required"`
	User    models.User  `gorm:"-" swaggerignore:"true"`
	Photo   models.Photo `gorm:"-" swaggerignore:"true"`
}

type CommentFormatter struct {
	ID      uint   `json:"id"`
	UserID  uint   `json:"user_id"`
	PhotoID uint   `json:"photo_id"`
	Message string `json:"message"`
}

func FormateComment(comment models.Comment) CommentFormatter {
	commentFormatter := CommentFormatter{
		ID:      comment.ID,
		UserID:  comment.UserID,
		PhotoID: comment.PhotoID,
		Message: comment.Message,
	}

	return commentFormatter
}

type CommentFormateDetail struct {
	ID      uint   `json:"id"`
	Message string `json:"message"`
	UserID  uint   `json:"user_id"`
	User    CommentUserFormatter
	PhotoID uint `json:"photo_id"`
	Photo   CommentPhotoFormatter
}

type CommentPhotoFormatter struct {
	ID      uint   `json:"photo_id"`
	Title   string `json:"title"`
	Caption string `json:"caption"`
}

type CommentUserFormatter struct {
	Username string `json:"username"`
	Age      int    `json:"age"`
}

func FormateCommentDetail(comment models.Comment) CommentFormateDetail {
	commentFormateDetail := CommentFormateDetail{
		ID:      comment.ID,
		Message: comment.Message,
		UserID:  comment.UserID,
		PhotoID: comment.PhotoID,
	}
	user := comment.User
	commentUserFormatter := CommentUserFormatter{
		Username: user.Username,
		Age:      user.Age,
	}
	commentFormateDetail.User = commentUserFormatter
	photo := comment.Photo
	commentPhotoFormatter := CommentPhotoFormatter{
		ID:      photo.ID,
		Title:   photo.Title,
		Caption: photo.Caption,
	}
	commentFormateDetail.Photo = commentPhotoFormatter

	return commentFormateDetail
}

func FormateCommentDetails(comments []models.Comment) []CommentFormateDetail {
	commentsFormatter := []CommentFormateDetail{}

	for _, comment := range comments {
		commentFormatter := FormateCommentDetail(comment)
		commentsFormatter = append(commentsFormatter, commentFormatter)
	}

	return commentsFormatter
}
