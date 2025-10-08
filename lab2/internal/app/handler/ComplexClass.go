package handler

import (
	"lab2/internal/app/ds"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) GetComplexClasses(ctx *gin.Context) {
	var ComplexClasses []ds.ComplexClass
	var err error

	searchQuery := ctx.Query("search_degree")
	if searchQuery == "" {
		ComplexClasses, err = h.Repository.GetComplexClasses()
		if err != nil {
			logrus.Error(err)
		}
	} else {
		ComplexClasses, err = h.Repository.GetComplexClasssByDegree(searchQuery)
		if err != nil {
			logrus.Error(err)
		}
	}

	ctx.HTML(http.StatusOK, "ComplexClasses.html", gin.H{
		"ComplexClasses":   ComplexClasses,
		"BigORequestId":    h.Repository.GetActiveBigORequestID(),
		"BigORequestCount": h.Repository.GetBigORequestCount(),
		"search_degree":    searchQuery,
	})
}

func (h *Handler) GetComplexClass(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	ComplexClass, err := h.Repository.GetComplexClass(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		logrus.Error(err)
		return
	}

	ctx.HTML(http.StatusOK, "ComplexClass.html", gin.H{
		"ComplexClass": gin.H{
			"ID":          ComplexClass.ID,
			"Complexity":  ComplexClass.Complexity,
			"Degree":      ComplexClass.Degree,
			"DegreeText":  ComplexClass.DegreeText,
			"Description": ComplexClass.Description,
			"IMG":         ComplexClass.IMG,
		},
	})

}
func (h *Handler) GetBigORequest(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
	}

	isDraft, err := h.Repository.IsDraftBigORequest(id)
	if !isDraft || err != nil {
		ctx.Redirect(http.StatusSeeOther, "/")
		return
	}

	BigORequestItems, err := h.Repository.GetBigORequest(id)
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "BigORequest.html", gin.H{
		"BigORequest":     BigORequestItems,
		"BigO_Request_ID": id,
	})
}

func (h *Handler) AddToBigORequest(ctx *gin.Context) {
	complexclassIDStr := ctx.PostForm("complexclass_id")
	complexclassID, err := strconv.Atoi(complexclassIDStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	creatorID := uint(1)

	err = h.Repository.AddComplexClass(uint(complexclassID), creatorID)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Redirect(http.StatusSeeOther, ctx.Request.Referer())
}

func (h *Handler) DeleteBigORequest(ctx *gin.Context) {
	bigorequestIDStr := ctx.PostForm("big_o_request_id")
	bigorequestID, err := strconv.Atoi(bigorequestIDStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	err = h.Repository.DeleteBigORequest(uint(bigorequestID))
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Redirect(http.StatusSeeOther, "/")
}
