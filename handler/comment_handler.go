package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/onainadapdap1/dev/kode/my_gram/dtos"
	"github.com/onainadapdap1/dev/kode/my_gram/models"
	"github.com/onainadapdap1/dev/kode/my_gram/service"
	"github.com/onainadapdap1/dev/kode/my_gram/utils"
)

type CommentHandlerInterface interface {
	CreateComment(c *gin.Context)
	UpdateComment(c *gin.Context)
	FindCommentByID(c *gin.Context)
	FindAllComments(c *gin.Context)
	DeleteCommentByID(c *gin.Context)
	RestoreCommentByID(c *gin.Context)
}

type commentHandler struct {
	service service.CommentServiceInterface
}

func NewCommentHandler(service service.CommentServiceInterface) CommentHandlerInterface {
	return &commentHandler{service: service}
}

// Create Comment  godoc
// @Summary Create comment
// @Description Create a new comment with given message and photo id
// @Tags comments
// @Accept json
// @Produce json
// @Param inputData dtos.CreateCommentInput body dtos.CreateCommentInput{} true "create new comment"
// @Success 200 {object} dtos.CommentFormatter
// @Failure 400 {object} utils.Response
// @Security BearerAuth
// @Router /api/v1/comments/comment [post]
func (h *commentHandler) CreateComment(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(models.User)

	var inputData dtos.CreateCommentInput
	if err := c.ShouldBindJSON(&inputData); err != nil {
		response := utils.APIResponse(fmt.Sprintf("failed to bind json data : %v", err), http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	photo, err := h.service.FindPhotoByID(uint(inputData.PhotoID))
	if err != nil {
		response := utils.APIResponse("failed to get photo by id", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return

	}
	// comment := models.Comment{}
	createCommentInput := dtos.CreateCommentInput{
		Message: inputData.Message,
		PhotoID: photo.ID,
		User:    currentUser,
		Photo:   photo,
	}

	newComment, err := h.service.CreateComment(createCommentInput)
	if err != nil {
		log.Printf("failed to create comment: %v", err)
		response := utils.APIResponse(fmt.Sprintf("failed to create comment : %v", err), http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := utils.APIResponse("success to create product", http.StatusOK, "success", dtos.FormateComment(newComment))
	c.JSON(http.StatusOK, response)
}

// Update Comment  godoc
// @Summary Update comment
// @Description Update a new comment with given message and photo id
// @Tags comments
// @Accept json
// @Produce json
// @Param id path int true "Photo iD"
// @Param inputData dtos.UpdateCommentInput body dtos.UpdateCommentInput{} true "update comment"
// @Success 200 {object} dtos.CommentFormatter
// @Failure 400 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Security BearerAuth
// @Router /api/v1/comments/comment/{id} [put]
func (h *commentHandler) UpdateComment(c *gin.Context) {
	var inputID dtos.GetCommentDetailInput
	if err := c.ShouldBindUri(&inputID); err != nil {
		response := utils.APIResponse("faile to get comment id", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var inputData dtos.UpdateCommentInput
	if err := c.ShouldBindJSON(&inputData); err != nil {
		response := utils.APIResponse("faile to get comment input", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	comment, err := h.service.FindCommentByID(inputID.ID)
	if err != nil {
		response := utils.APIResponse("Failed find comment by id", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	currentUser := c.MustGet("currentUser").(models.User)
	userId := currentUser.ID
	inputData.User = currentUser
	// inputData.Photo = photo

	if userId != comment.UserID {
		response := utils.APIResponse("Unauthorized", http.StatusForbidden, "error", nil)
		c.JSON(http.StatusForbidden, response)
		return
	}

	updatedComment, err := h.service.UpdateComment(inputID, inputData)
	if err != nil {
		log.Printf("failed to update comment: %v", err)
		response := utils.APIResponse(fmt.Sprintf("failed to update comment: %v", err), http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := utils.APIResponse("success to update comment", http.StatusOK, "success", dtos.FormateComment(updatedComment))
	c.JSON(http.StatusOK, response)

}

// Get Comment by ID godoc
// @Summary Get one comment by id
// @Description Get one comment by id
// @Tags comments
// @Produce json
// @Param id path int true "get comment by id"
// @Success 200 {object} dtos.CommentFormateDetail{}
// @Failure 400 {object} utils.Response
// @Router /api/v1/comments/comment/{id} [get]
func (h *commentHandler) FindCommentByID(c *gin.Context) {
	var inputID dtos.GetCommentDetailInput

	if err := c.ShouldBindUri(&inputID); err != nil {
		response := utils.APIResponse("faile to get comment id", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	comment, err := h.service.FindCommentByID(inputID.ID)
	if err != nil {
		response := utils.APIResponse("Failed find comment by id", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := utils.APIResponse("success to get product detail", http.StatusOK, "success", dtos.FormateCommentDetail(comment))
	c.JSON(http.StatusOK, response)
}

// Get all comments godoc
// @Summary Get all comments
// @Description Get all comments
// @Tags comments
// @Produce json
// @Success 200 {object} []dtos.CommentFormateDetail{}
// @Failure 400 {object} utils.Response
// @Router /api/v1/comments [get]
func (h *commentHandler) FindAllComments(c *gin.Context) {
	comments, err := h.service.FindAllComments()
	if err != nil {
		response := utils.APIResponse("Failed get all comments", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := utils.APIResponse("success to get all comments", http.StatusOK, "success", dtos.FormateCommentDetails(comments))
	c.JSON(http.StatusOK, response)
}

// Delete comment by ID godoc
// @Summary Delete comment by id
// @Description Delete comment by id
// @Tags comments
// @Produce json
// @Param id path int true "delete comment by id"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Security BearerAuth
// @Router /api/v1/comments/comment/{id} [delete]
func (h *commentHandler) DeleteCommentByID(c *gin.Context) {
	var inputID dtos.GetCommentDetailInput
	if err := c.ShouldBindUri(&inputID); err != nil {
		response := utils.APIResponse("faile to get comment id", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	currentUser := c.MustGet("currentUser").(models.User)
	userId := currentUser.ID

	comment, err := h.service.FindCommentByID(uint(inputID.ID))
	if err != nil {
		response := utils.APIResponse("Failed to find comment by id", http.StatusInternalServerError, "error", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	if userId != comment.UserID {
		response := utils.APIResponse("Unauthorized", http.StatusForbidden, "error", nil)
		c.JSON(http.StatusForbidden, response)
		return
	}

	if err := h.service.DeleteCommentByID(uint(inputID.ID)); err != nil {
		response := utils.APIResponse("Failed to delete comment", http.StatusInternalServerError, "error", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := utils.APIResponse("success to delete comment", http.StatusOK, "success", nil)
	c.JSON(http.StatusOK, response)
}

// Restore comment by ID godoc
// @Summary Restore comment by id
// @Description Restore comment by id
// @Tags comments
// @Param id path int true "restore comment by id"
// @Success 200 {object} dtos.CommentFormateDetail
// @Failure 400 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Security BearerAuth
// @Router /api/v1/comments/restorecomment/{id} [put]
func (h *commentHandler) RestoreCommentByID(c *gin.Context) {
	var inputID dtos.GetCommentDetailInput
	if err := c.ShouldBindUri(&inputID); err != nil {
		response := utils.APIResponse("faile to get comment id", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	currentUser := c.MustGet("currentUser").(models.User)
	userId := currentUser.ID

	comment, err := h.service.GetCommentByID(uint(inputID.ID))
	if err != nil {
		response := utils.APIResponse("Failed to find comment by id", http.StatusInternalServerError, "error", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	if userId != comment.UserID {
		response := utils.APIResponse("Unauthorized", http.StatusForbidden, "error", nil)
		c.JSON(http.StatusForbidden, response)
		return
	}

	restoredComment, err := h.service.RestoreCommentByID(uint(inputID.ID))
	if err != nil {
		response := utils.APIResponse("Failed to restore comment", http.StatusInternalServerError, "error", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := utils.APIResponse("success to restore comment", http.StatusOK, "success", dtos.FormateCommentDetail(restoredComment))
	c.JSON(http.StatusOK, response)
}
