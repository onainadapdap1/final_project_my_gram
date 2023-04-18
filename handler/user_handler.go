package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/onainadapdap1/dev/kode/my_gram/dtos"
	"github.com/onainadapdap1/dev/kode/my_gram/helpers"
	"github.com/onainadapdap1/dev/kode/my_gram/models"
	"github.com/onainadapdap1/dev/kode/my_gram/service"
	"github.com/onainadapdap1/dev/kode/my_gram/utils"
)

type UserHandlerInterface interface {
	RegisterUser(c *gin.Context)
	LoginUser(c *gin.Context)
}

type userHandler struct {
	service service.UserServiceInterface
}

func NewUserHandler(service service.UserServiceInterface) UserHandlerInterface {
	return &userHandler{service: service}
}

func (h *userHandler) RegisterUser(c *gin.Context) {
	var input models.RegisterUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		errors := utils.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := utils.APIResponse(fmt.Sprintf("%v", err), http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	newUser, err := h.service.Registeruser(input)
	if err != nil {
		response := utils.APIResponse(fmt.Sprintf("%v", err), http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	userFormateReg := dtos.FormateUserRegister(newUser)
	response := utils.APIResponse("Account has been registered", http.StatusCreated, "success", userFormateReg)
	c.JSON(http.StatusOK, response)
}

func (h *userHandler) LoginUser(c *gin.Context) {
	var input models.LoginUserInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := utils.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := utils.APIResponse("failed to bind data", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	loggeinUser, err := h.service.LoginUser(input)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}

		response := utils.APIResponse("Login failed user data", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	token, err := helpers.GenerateToken(loggeinUser.ID, loggeinUser.Email)
	if err != nil {
		response := utils.APIResponse("Login failed to generate token", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	loginFormateUser := dtos.FormateUserLogin(loggeinUser, token)
	
	response := utils.APIResponse("Successfully loggedin", http.StatusOK, "success", loginFormateUser)
	c.JSON(http.StatusOK, response)
}