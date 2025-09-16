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

	searchQuery := ctx.Query("query")
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
	cartComplexClasses, err := h.Repository.GetCart()
	cartCount := 0
	if err == nil {
		cartCount = len(cartComplexClasses)
	}
	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"complexityClasses": complexityClasses,
		"query":             searchQuery,
		"cartCount":         cartCount,
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

	ctx.HTML(http.StatusOK, "order.html", gin.H{
		"complexClass": complexClass,
	})
}

func (h *Handler) GetCart(ctx *gin.Context) {
	var complexityClasses []repository.ComplexClass
	var err error

	complexityClasses, err = h.Repository.GetCart()
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "cart.html", gin.H{
		"service_complexclasses": complexityClasses,
	})
}
