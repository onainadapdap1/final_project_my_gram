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

type SocialmediaHandlerInterface interface {
	CreateSocialMedia(c *gin.Context)
	UpdateSocialMedia(c *gin.Context)
	FindAllSocialMedia(c *gin.Context)
	FindBySocialMediaID(c *gin.Context)
	DeleteSocialMedia(c *gin.Context)
	RestoreSocialMedia(c *gin.Context)
}

type socialmediaHandler struct {
	service service.SocialmediaServiceInterface
}

func NewSocialmediaHandler(service service.SocialmediaServiceInterface) SocialmediaHandlerInterface {
	return &socialmediaHandler{service: service}
}

// Create Social Media godoc
// @Summary Create a new social media
// @Description Create new social media
// @Tags socialmedias
// @Accept json
// @Produce json
// @Param input dtos.CreateSocialMediaInput body dtos.CreateSocialMediaInput{} true "create user social media"
// @Success 200 {object} dtos.SocialMediaFormatter
// @Failure 400 {object} utils.Response
// @Security BearerAuth
// @Router /api/v1/socialmedias/socialmedia [post]
func (h *socialmediaHandler) CreateSocialMedia(c *gin.Context) {
	var input dtos.CreateSocialMediaInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response := utils.APIResponse(fmt.Sprintf("%v", err), http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// get current user
	currentUser := c.MustGet("currentUser").(models.User)

	input.User = currentUser
	newSocialMedia, err := h.service.CreateSocialMedia(&input)
	if err != nil {
		log.Printf("failed to create social media: %v", err)
		response := utils.APIResponse(fmt.Sprintf("%v", err), http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := utils.APIResponse("Success to create social media", http.StatusOK, "success", dtos.FormateSocialMedia(newSocialMedia))
	c.JSON(http.StatusOK, response)
}

// Update Social Media godoc
// @Summary Update social media by id
// @Description Update social media by id
// @Tags socialmedias
// @Accept json
// @Produce json
// @Param id path int true "social media iD"
// @Param inputData dtos.CreateSocialMediaInput body dtos.CreateSocialMediaInput{} true "update comment"
// @Success 200 {object} dtos.SocialMediaFormatter
// @Failure 400 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Security BearerAuth
// @Router /api/v1/socialmedias/socialmedia/{id} [put]
func (h *socialmediaHandler) UpdateSocialMedia(c *gin.Context) {
	var inputID dtos.GetSocialMediaDetailInput
	if err := c.ShouldBindUri(&inputID); err != nil {
		response := utils.APIResponse(fmt.Sprintf("Failed to get social media id : %v", err), http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var inputData dtos.CreateSocialMediaInput
	if err := c.ShouldBindJSON(&inputData); err != nil {
		response := utils.APIResponse(fmt.Sprintf("%v", err), http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	socialMedia, err := h.service.FindSocialMediaByID(&inputID)
	if err != nil {
		response := utils.APIResponse(fmt.Sprintf("Failed find by id : %v", err), http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	currentUser := c.MustGet("currentUser").(models.User)
	userId := currentUser.ID

	if userId != socialMedia.UserID {
		response := utils.APIResponse("Unauthorized", http.StatusForbidden, "error", nil)
		c.JSON(http.StatusForbidden, response)
		return
	}

	updatedSocialMedia, err := h.service.UpdateSocialMedia(&inputID, &inputData)
	if err != nil {
		log.Printf("failed to update social media : %v", err)
		response := utils.APIResponse(fmt.Sprintf("failed to update social media: %v", err), http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := utils.APIResponse("sucess to update social media", http.StatusOK, "success", dtos.FormateSocialMediaDetail(updatedSocialMedia))
	c.JSON(http.StatusOK, response)

}

// Get all social medias godoc
// @Summary Get all social medias
// @Description Get all social medias
// @Tags socialmedias
// @Produce json
// @Success 200 {object} []dtos.SocialMediaDetailFormatter{}
// @Failure 400 {object} utils.Response
// @Router /api/v1/socialmedias [get]
func (h *socialmediaHandler) FindAllSocialMedia(c *gin.Context) {
	socialMedias, err := h.service.FindAllSocialMedia()
	if err != nil {
		response := utils.APIResponse(fmt.Sprintf("Failed find by id : %v", err), http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := utils.APIResponse("list of social medias", http.StatusOK, "success", dtos.FormateSocialMediaDetails(socialMedias))
	c.JSON(http.StatusOK, response)
}
// Get Social Media by ID godoc
// @Summary Get one social media by id
// @Description Get one social media by id
// @Tags socialmedias
// @Produce json
// @Param id path int true "get social media by id"
// @Success 200 {object} dtos.SocialMediaDetailFormatter{}
// @Failure 400 {object} utils.Response
// @Router /api/v1/socialmedias/socialmedia/{id} [get]
func (h *socialmediaHandler) FindBySocialMediaID(c *gin.Context) {
	var inputID dtos.GetSocialMediaDetailInput
	if err := c.ShouldBindUri(&inputID); err != nil {
		response := utils.APIResponse(fmt.Sprintf("Failed to get social media id : %v", err), http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	socialMedia, err := h.service.FindSocialMediaByID(&inputID)
	if err != nil {
		response := utils.APIResponse(fmt.Sprintf("Failed find by id : %v", err), http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := utils.APIResponse("sucess to get social media", http.StatusOK, "success", dtos.FormateSocialMediaDetail(socialMedia))
	c.JSON(http.StatusOK, response)
}

// Delete social media by ID godoc
// @Summary Delete social media by id
// @Description Delete social media by id
// @Tags socialmedias
// @Produce json
// @Param id path int true "delete social media by id"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Security BearerAuth
// @Router /api/v1/socialmedias/sosmed/{id} [delete]
func (h *socialmediaHandler) DeleteSocialMedia(c *gin.Context) {
	var inputID dtos.GetSocialMediaDetailInput
	if err := c.ShouldBindUri(&inputID); err != nil {
		response := utils.APIResponse(fmt.Sprintf("Failed to get social media id : %v", err), http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	socialMedia, err := h.service.FindSocialMediaByID(&inputID)
	if err != nil {
		response := utils.APIResponse(fmt.Sprintf("Failed find by id : %v", err), http.StatusUnprocessableEntity, "error", nil)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := c.MustGet("currentUser").(models.User)
	userId := currentUser.ID

	if userId != socialMedia.UserID {
		response := utils.APIResponse("Unauthorized", http.StatusForbidden, "error", nil)
		c.JSON(http.StatusForbidden, response)
		return
	}

	err = h.service.DeleteSocialMedia(inputID.ID)
	if err != nil {
		response := utils.APIResponse(fmt.Sprintf("Failed to delete social media : %v", err), http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := utils.APIResponse("sucess to delete social media", http.StatusOK, "success", nil)
	c.JSON(http.StatusOK, response)
}

// Restore social media by ID godoc
// @Summary Restore social media by id
// @Description Restore social media by id
// @Tags socialmedias
// @Param id path int true "restore social media by id"
// @Success 200 {object} dtos.SocialMediaDetailFormatter
// @Failure 400 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Security BearerAuth
// @Router /api/v1/socialmedias/restoresosmed/{id} [put]
func (h *socialmediaHandler) RestoreSocialMedia(c *gin.Context) {
	var inputID dtos.GetSocialMediaDetailInput
	if err := c.ShouldBindUri(&inputID); err != nil {
		response := utils.APIResponse(fmt.Sprintf("Failed to get social media id : %v", err), http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	socialMedia, err := h.service.GetDeletedSocialMediaByID(uint(inputID.ID))
	if err != nil {
		response := utils.APIResponse(fmt.Sprintf("Failed find by id : %v", err), http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	currentUser := c.MustGet("currentUser").(models.User)
	userId := currentUser.ID

	if userId != socialMedia.UserID {
		response := utils.APIResponse("Unauthorized", http.StatusForbidden, "error", nil)
		c.JSON(http.StatusForbidden, response)
		return
	}

	restoredSocialMedia, err := h.service.RestoreSocialMedia(inputID.ID)
	if err != nil {
		response := utils.APIResponse(fmt.Sprintf("Failed to restore social media : %v", err), http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := utils.APIResponse("sucess to restore social media", http.StatusOK, "success", dtos.FormateSocialMediaDetail(restoredSocialMedia))
	c.JSON(http.StatusOK, response)
}
