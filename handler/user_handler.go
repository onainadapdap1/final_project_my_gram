package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/onainadapdap1/dev/kode/my_gram/dtos"
	"github.com/onainadapdap1/dev/kode/my_gram/helpers"
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

// Register User godoc
// @Summary Register a new user
// @Description Register new user
// @Tags users
// @Accept json
// @Produce json
// @Param input dtos.RegisterUserInput body dtos.RegisterUserInput{} true "register user"
// @Success 200 {object} dtos.UserRegisterFormatter
// @Failure 400 {object} utils.Response
// @Failure 422 {object} utils.Response
// @Router /api/v1/register [post]
func (h *userHandler) RegisterUser(c *gin.Context) {
	var input dtos.RegisterUserInput
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

// Login User godoc
// @Summary Login user
// @Description User login with email and password
// @Tags users
// @Accept json
// @Produce json
// @Param input dtos.LoginUserInput body dtos.LoginUserInput{} true "Login user input"
// @Success 200 {object} dtos.UserRegisterFormatter
// @Failure 400 {object} utils.Response
// @Failure 422 {object} utils.Response
// @Router /api/v1/login [post]
func (h *userHandler) LoginUser(c *gin.Context) {
	var input dtos.LoginUserInput

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
