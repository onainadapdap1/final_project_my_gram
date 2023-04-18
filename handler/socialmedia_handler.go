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

	response := utils.APIResponse("sucess to update social media", http.StatusOK, "success", dtos.FormateSocialMedia(updatedSocialMedia))
	c.JSON(http.StatusOK, response)

}

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

func (h *socialmediaHandler) DeleteSocialMedia(c *gin.Context) {
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
