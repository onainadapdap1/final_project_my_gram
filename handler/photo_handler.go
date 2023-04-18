package handler

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/onainadapdap1/dev/kode/my_gram/dtos"
	"github.com/onainadapdap1/dev/kode/my_gram/models"
	"github.com/onainadapdap1/dev/kode/my_gram/service"
	"github.com/onainadapdap1/dev/kode/my_gram/utils"
)

type PhotoHandlerInterface interface {
	CreatePhoto(c *gin.Context)
	UpdatePhoto(c *gin.Context)
	FindAllPhoto(c *gin.Context)
	FindByPhotoID(c *gin.Context)
	DeletePhotoByID(c *gin.Context) 
}

type photoHandler struct {
	service service.PhotoServiceInterface
}

func NewPhotoHandler(service service.PhotoServiceInterface) PhotoHandlerInterface {
	return &photoHandler{service: service}
}

func (h *photoHandler) CreatePhoto(c *gin.Context) {
	title := c.PostForm("title")
	caption := c.PostForm("caption")
	file, err := c.FormFile("photo_url")
	if err != nil {
		// response := utils.APIResponse("Failed to bind image file", http.StatusBadRequest, "error", nil)
		response := utils.APIResponse(fmt.Sprintf("%v",err), http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	currentUser := c.MustGet("currentUser").(models.User)
	userId := currentUser.ID

	fileName := fmt.Sprintf("%d-%s", userId, file.Filename)
	dirPath := filepath.Join(".", "static", "images", "photos")
	filePath := filepath.Join(dirPath, fileName)

	// Create directory if does not exist
	if _, err = os.Stat(dirPath); os.IsNotExist(err) {
		err = os.MkdirAll(dirPath, 0755)
		if err != nil {
			response := utils.APIResponse("Failed to create directory", http.StatusBadRequest, "error", nil)
			c.JSON(http.StatusInternalServerError, response)
			return
		}
	}

	// Create file that will hold the image
	outputFile, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

	// Open the temporary file that contains the uploaded image
	inputFile, err := file.Open()
	if err != nil {
		response := utils.APIResponse("Failed to open temporary photo file", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusOK, response)
	}
	defer inputFile.Close()

	// Copy the temporary image to the permanent location outputFile
	_, err = io.Copy(outputFile, inputFile)
	if err != nil {
		log.Fatal(err)
		c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
		return
	}

	createPhotoInput := dtos.CreatePhotoInput{
		User:     currentUser,
		Title:    title,
		Caption:  caption,
		PhotoUrl: filePath,
	}

	newPhoto, err := h.service.CreatePhoto(createPhotoInput)
	if err != nil {
		log.Printf("failed to create photo: %v", err)
		// response := utils.APIResponse("Failed to create photo", http.StatusBadRequest, "error", nil)
		response := utils.APIResponse(fmt.Sprintf("%v", err), http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := utils.APIResponse("Success to create photo", http.StatusOK, "success", dtos.FormatePhoto(newPhoto))
	c.JSON(http.StatusOK, response)
}

func (h *photoHandler) UpdatePhoto(c *gin.Context) {
	var inputID dtos.GetPhotoDetailInput
	err := c.ShouldBindUri(&inputID)
	if err != nil {
		response := utils.APIResponse(fmt.Sprintf("Failed to get photo id : %v", err), http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var inputData dtos.CreatePhotoInput
	title := c.PostForm("title")
	caption := c.PostForm("caption")

	photo, err := h.service.FindByPhotoID(inputID.ID)
	if err != nil {
		response := utils.APIResponse(fmt.Sprintf("Failed find by id : %v", err), http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	currentUser := c.MustGet("currentUser").(models.User)
	userId := currentUser.ID

	if userId != photo.UserID {
		response := utils.APIResponse("Unauthorized", http.StatusForbidden, "error", nil)
		c.JSON(http.StatusForbidden, response)
		return
	}

	// handle file after get data
	file, err := c.FormFile("photo_url")
	if err != nil {
		// use existing image url if file is not found
		inputData.PhotoUrl = photo.PhotoUrl
		response := utils.APIResponse("Failed to upload file", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	} else {
		// remove the old image file from the static folder,
		if photo.PhotoUrl != "" {
			// oldFilename := strings.TrimPrefix(category.ImageURL, "/static/images/categories/")
			oldFilename := filepath.Base(photo.PhotoUrl)
			if err := os.Remove("static/images/photos/" + oldFilename); err != nil {
				log.Printf("Failed to remove old filename: %v", err)
				response := utils.APIResponse(fmt.Sprintf("Failed to remove old filename: %v", err), http.StatusInternalServerError, "error", nil)
				c.JSON(http.StatusInternalServerError, response)
				return
			}
		}

		fileName := fmt.Sprintf("%d-%s", userId, file.Filename)

		dirPath := filepath.Join(".", "static", "images", "photos")
		filePath := filepath.Join(dirPath, fileName)
		// Create directory if does not exist
		if _, err = os.Stat(dirPath); os.IsNotExist(err) {
			err = os.MkdirAll(dirPath, 0755)
			if err != nil {
				response := utils.APIResponse("Failed to upload photo image", http.StatusBadRequest, "error", nil)
				c.JSON(http.StatusInternalServerError, response)
				return
			}
		}
		// Create file that will hold the image
		outputFile, err := os.Create(filePath)
		if err != nil {
			log.Fatal(err)
		}
		defer outputFile.Close()

		// Open the temporary file that contains the uploaded image
		inputFile, err := file.Open()
		if err != nil {
			response := utils.APIResponse("Failed to upload photo image", http.StatusBadRequest, "error", nil)
			c.JSON(http.StatusOK, response)
		}
		defer inputFile.Close()

		// Copy the temporary image to the permanent location outputFile
		_, err = io.Copy(outputFile, inputFile)
		if err != nil {
			log.Fatal(err)
			c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
			return
		}

		inputData.PhotoUrl = filePath
	}
	inputData.Title = title
	inputData.Caption = caption
	updatedPhoto, err := h.service.UpdatePhoto(inputID, inputData)

	if err != nil {
		log.Printf("failed to update photo : %v", err)
		response := utils.APIResponse(fmt.Sprintf("failed to update photo: %v", err), http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := utils.APIResponse("sucess to update category", http.StatusOK, "success", dtos.FormatePhoto(updatedPhoto))
	c.JSON(http.StatusOK, response)

}

func (h *photoHandler) FindAllPhoto(c *gin.Context) {
	photos, err := h.service.FindAllPhoto()
	if err != nil {
		response := utils.APIResponse("failed to get all photo", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := utils.APIResponse("list of photos", http.StatusOK, "success", dtos.FormatePhotoDetails(photos))
	c.JSON(http.StatusOK, response)

}

func (h *photoHandler) FindByPhotoID(c *gin.Context) {
	param := c.Param("id")
	photoID, err := strconv.Atoi(param)
	// err := c.ShouldBindUri(&inputID)
	if err != nil {
		response := utils.APIResponse("failed to get detail photo input", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	photo, err := h.service.FindByPhotoID(uint(photoID))
	if err != nil {
		response := utils.APIResponse("failed to get detail photo", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := utils.APIResponse("success get photo detail", http.StatusOK, "success", dtos.FormatePhotoDetail(photo))
	c.JSON(http.StatusOK, response)
}

func (h *photoHandler) DeletePhotoByID(c *gin.Context) {
	param := c.Param("id")
	photoID, err := strconv.Atoi(param)
	if err != nil {
		response := utils.APIResponse("failed to get detail input", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	photo, err := h.service.FindByPhotoID(uint(photoID))
	if err != nil {
		response := utils.APIResponse("failed to get detail photo", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	err = h.service.DeletePhotoByID(photo)
	if err != nil {
		response := utils.APIResponse("failed to delete photo", http.StatusInternalServerError, "error", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := utils.APIResponse("Success to delete photo", http.StatusOK, "success", nil)
	c.JSON(http.StatusOK, response)

}