package handler

import (
	"lab1/internal/app/repository"
	"net/http"
	"strconv"

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

func (h *Handler) GetComplexClasses(ctx *gin.Context) {
	var complexityClasses []repository.ComplexClass
	var err error

	searchQuery := ctx.Query("search_degree")
	if searchQuery == "" {
		complexityClasses, err = h.Repository.GetComplexClasses()
		if err != nil {
			logrus.Error(err)
		}
	} else {
		complexityClasses, err = h.Repository.GetComplexClassByDegree(searchQuery)
		if err != nil {
			logrus.Error(err)
		}
	}

	RequestComplexClasses, err := h.Repository.GetBigORequest()
	RequestCount := 0
	if err == nil {
		RequestCount = len(RequestComplexClasses)
	}
	ctx.HTML(http.StatusOK, "ComplexClasses.html", gin.H{
		"complexityClasses": complexityClasses,
		"search_degree":     searchQuery,
		"BigORequestCount":  RequestCount,
	})
}
func (h *Handler) GetComplexClass(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
	}

	complexClass, err := h.Repository.GetComplexClass(id)
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "ComplexClass.html", gin.H{
		"complexClass": complexClass,
	})
}

func (h *Handler) GetBigORequest(ctx *gin.Context) {
	var complexityClasses []repository.ComplexClass
	var err error

	complexityClasses, err = h.Repository.GetBigORequest()
	if err != nil {
		logrus.Error(err)
	}
	resultTime := "13 мкс"
	resultComplexity := "O(n)"
	ctx.HTML(http.StatusOK, "BigORequest.html", gin.H{
		"service_complexclasses": complexityClasses,
		"ResultTime":             resultTime,
		"ResultComplexity":       resultComplexity,
	})
}
