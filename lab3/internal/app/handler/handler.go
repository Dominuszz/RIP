package handler

import (
	"errors"
	"lab3/internal/app/repository"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

type Handler struct {
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{
		Repository: r,
	}
}

// RegisterHandler godoc
// @title Big O Request API
// @version 1.0
// @description API для управления расчётами времени и сложности Классов сложности
// @contact.name API Support
// @contact.url http://localhost:8080
// @contact.email support@bigorequest.com
// @license.name MIT
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func (h *Handler) RegisterHandler(router *gin.Engine) {
	api := router.Group("/api/v1")

	unauthorized := api.Group("/")
	unauthorized.POST("/users/signup", h.CreateUser)
	unauthorized.POST("/users/signin", h.SignIn)
	unauthorized.GET("/complexclass", h.GetComplexClasses)
	unauthorized.GET("/complexclass/:id", h.GetComplexClass)

	authorized := api.Group("/")
	authorized.Use(h.ModeratorMiddleware(false))
	authorized.POST("/complexclass/create-compclass", h.CreateComplexClass)
	authorized.PUT("/complexclass/:id/edit-compclass", h.EditComplexClass)
	authorized.DELETE("/complexclass/:id/delete-compclass", h.DeleteComplexClass)
	authorized.POST("/complexclass/:id/add-to-bigorequest", h.AddToBigORequest)
	authorized.POST("/complexclass/:id/add-photo", h.AddPhoto)

	authorized.GET("/bigorequest/bigorequest-cart", h.GetBigORequestCart)
	authorized.GET("/bigorequest/all-bigo_requests", h.GetAllBigORequests)
	authorized.GET("/bigorequest/:id", h.GetBigORequest)
	authorized.PUT("/bigorequest/:id/edit-bigorequest", h.EditBigORequest)
	authorized.PUT("/bigorequest/:id/form-bigorequest", h.FormBigORequest)
	authorized.PUT("/bigorequest/:id/finish-bigorequest", h.FinishBigORequest)
	authorized.DELETE("/bigorequest/:id/delete-bigorequest", h.DeleteBigORequest)

	authorized.DELETE("/compclassrequest/:compclass_id/:bigo_request_id", h.DeleteCompClassFromBigORequest)
	authorized.PUT("/compclassrequest/:compclass_id/:bigo_request_id", h.EditCompClassFromBigORequest)

	authorized.GET("/users/:login/info", h.GetInfo)
	authorized.PUT("/users/:login/info", h.EditInfo)
	authorized.POST("/users/signout", h.SignOut)

	moderator := api.Group("/")
	moderator.Use(h.ModeratorMiddleware(true))
	authorized.PUT("/bigorequest/:id/form", h.FormBigORequest)

	swaggerURL := ginSwagger.URL("/swagger/doc.json")
	router.Any("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, swaggerURL))
	router.GET("/swagger", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})
}

func (h *Handler) RegisterStatic(router *gin.Engine) {
	router.LoadHTMLGlob("templates/*")
	router.Static("/styles", "./resources/styles")
	router.Static("/img", "./resources/img")
}

func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) {
	logrus.Error(err.Error())

	var errorMessage string
	switch {
	case errors.Is(err, repository.ErrNotFound):
		errorMessage = "Не найден"
	case errors.Is(err, repository.ErrAlreadyExists):
		errorMessage = "Уже существует"
	case errors.Is(err, repository.ErrNotAllowed):
		errorMessage = "Доступ запрещен"
	case errors.Is(err, repository.ErrNoDraft):
		errorMessage = "Черновик не найден"
	default:
		errorMessage = err.Error()
	}

	ctx.JSON(errorStatusCode, gin.H{
		"status":      "error",
		"description": errorMessage,
	})
}
