package handler

import (
	"errors"
	"lab3/internal/app/repository"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{
		Repository: r,
	}
}

func (h *Handler) RegisterHandler(router *gin.Engine) {
	router.GET("/api/complexclass", h.GetComplexClasses)
	router.GET("/api/complexclass/:id", h.GetComplexClass)
	router.POST("/api/complexclass/create-compclass", h.CreateComplexClass)
	router.PUT("/api/complexclass/:id/edit-compclass", h.EditComplexClass)
	router.DELETE("/api/complexclass/:id/delete-compclass", h.DeleteComplexClass)
	router.POST("/api/complexclass/:id/add-to-bigorequest", h.AddToBigORequest)
	router.POST("/api/complexclass/:id/add-photo", h.AddPhoto)

	router.GET("/api/bigorequest/bigorequest-cart", h.GetBigORequestCart)
	router.GET("/api/bigorequest/all-bigo_requests", h.GetAllBigORequests)
	router.GET("/api/bigorequest/:id", h.GetBigORequest)
	router.PUT("/api/bigorequest/:id/edit-bigorequest", h.EditBigORequest)
	router.PUT("/api/bigorequest/:id/form-bigorequest", h.FormBigORequest)
	router.PUT("/api/bigorequest/:id/finish-bigorequest", h.FinishBigORequest)
	router.DELETE("/api/bigorequest/:id/delete-bigorequest", h.DeleteBigORequest)

	router.DELETE("/api/compclassrequest/:compclass_id/:bigo_request_id", h.DeleteCompClassFromBigORequest)
	router.PUT("/api/compclassrequest/:compclass_id/:bigo_request_id", h.EditCompClassFromBigORequest)

	router.POST("/api/users/signup", h.CreateUser)
	router.GET("/api/users/info", h.GetInfo)
	router.PUT("/api/users/info", h.EditInfo)
	router.POST("/api/users/signin", h.SignIn)
	router.POST("/api/users/signout", h.SignOut)
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
