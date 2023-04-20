package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/onainadapdap1/dev/kode/my_gram/driver"
	"github.com/onainadapdap1/dev/kode/my_gram/handler"
	"github.com/onainadapdap1/dev/kode/my_gram/middlewares"
	"github.com/onainadapdap1/dev/kode/my_gram/repository"
	"github.com/onainadapdap1/dev/kode/my_gram/service"

	_ "github.com/onainadapdap1/dev/kode/my_gram/docs"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title 			MyGram API
// @version         1.0
// @description     This is service for MyGram.
// @termsOfService  http://swagger.io/terms/

// @contact.name    API Support
// @contact.url     http://www.swagger.io/support
// @contact.email   support@swagger.io

// @license.name  	Apache 2.0
// @license.url   	http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func Router() *gin.Engine {
	router := gin.Default()

	db := driver.ConnectDB()
	
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	photoRepo := repository.NewPhotoRepository(db)
	photoService := service.NewPhotoService(photoRepo)
	photoHandler := handler.NewPhotoHandler(photoService)

	commentRepo := repository.NewCommentRepository(db)
	commentService := service.NewCommentService(commentRepo)
	commentHandler := handler.NewCommentHandler(commentService)

	socialMediaRepo := repository.NewSocialmediaRepository(db)
	socialMediaService := service.NewSocialmediaService(socialMediaRepo)
	socialMediaHandler := handler.NewSocialmediaHandler(socialMediaService)

	api := router.Group("api/v1")
	api.POST("/register", userHandler.RegisterUser)
	api.POST("/login", userHandler.LoginUser)

	photoRouter := api.Group("/photos")
	{
		photoRouter.GET("", photoHandler.FindAllPhoto)
		photoRouter.GET("/photo/:id", photoHandler.FindByPhotoID)
		photoRouter.Use(middlewares.Authentication())
		photoRouter.POST("/photo", userAuthorization(userService), photoHandler.CreatePhoto)
		photoRouter.PUT("/photo/:id", userAuthorization(userService), photoHandler.UpdatePhoto)
		photoRouter.DELETE("/photo/:id", userAuthorization(userService), photoHandler.DeletePhotoByID)
	}
	commentRouter := api.Group("/comments")
	{
		commentRouter.GET("", commentHandler.FindAllComments)
		commentRouter.GET("/comment/:id", commentHandler.FindCommentByID)
		commentRouter.Use(middlewares.Authentication())
		commentRouter.POST("/comment", userAuthorization(userService), commentHandler.CreateComment)
		commentRouter.PUT("/comment/:id", userAuthorization(userService), commentHandler.UpdateComment)
		commentRouter.DELETE("/comment/:id", userAuthorization(userService), commentHandler.DeleteCommentByID)
		commentRouter.PUT("/restorecomment/:id", userAuthorization(userService), commentHandler.RestoreCommentByID)
	}
	socialmediaRouter := api.Group("/socialmedias")
	{
		socialmediaRouter.GET("", socialMediaHandler.FindAllSocialMedia)
		socialmediaRouter.GET("/socialmedia/:id", socialMediaHandler.FindBySocialMediaID)
		socialmediaRouter.Use(middlewares.Authentication())
		socialmediaRouter.POST("/socialmedia", userAuthorization(userService), socialMediaHandler.CreateSocialMedia)
		socialmediaRouter.PUT("/socialmedia/:id", userAuthorization(userService), socialMediaHandler.UpdateSocialMedia)
		socialmediaRouter.DELETE("/sosmed/:id", userAuthorization(userService), socialMediaHandler.DeleteSocialMedia)
		socialmediaRouter.PUT("/restoresosmed/:id", userAuthorization(userService), socialMediaHandler.RestoreSocialMedia)
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	return router
}

func userAuthorization(userSerivce service.UserServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		userData := c.MustGet("userData").(jwt.MapClaims)
		userId := uint(userData["user_id"].(float64))

		user, err := userSerivce.GetUserByID(userId)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error":   "Data not found",
				"message": "data doesn't exist",
			})
			return
		}

		if user.ID == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "you are now allowed to access this data",
			})
			return
		}

		c.Set("currentUser", user)
		c.Next()
	}
}
