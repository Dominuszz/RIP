package handler

import (
	"errors"
	"fmt"
	"lab3/internal/app/ds"
	"lab3/internal/app/repository"
	"lab3/internal/app/serializer"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetComplexClasses(ctx *gin.Context) {
	var compclasses []ds.ComplexClass
	var err error

	searchQuery := ctx.Query("search-degree")
	if searchQuery == "" {
		compclasses, err = h.Repository.GetComplexClasses()
		if err != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
			return
		}
	} else {
		compclasses, err = h.Repository.GetCompClassByDegree(searchQuery)
		if err != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
			return
		}
	}
	resp := make([]serializer.ComplexClassJSON, 0, len(compclasses))
	for _, r := range compclasses {
		resp = append(resp, serializer.CompClassToJSON(r))
	}
	ctx.JSON(http.StatusOK, resp)
}

func (h *Handler) GetComplexClass(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	compclass, err := h.Repository.GetComplexClass(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, serializer.CompClassToJSON(*compclass))
}

func (h *Handler) CreateComplexClass(ctx *gin.Context) {
	var compclassJSON serializer.ComplexClassJSON
	if err := ctx.BindJSON(&compclassJSON); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	compclass, err := h.Repository.CreateComplexClass(compclassJSON)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Header("Location", fmt.Sprintf("/ComplexClass/%v", compclass.ID))
	ctx.JSON(http.StatusCreated, serializer.CompClassToJSON(compclass))
}

func (h *Handler) DeleteComplexClass(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	err = h.Repository.DeleteComplexClass(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "deleted",
	})
}

func (h *Handler) EditComplexClass(ctx *gin.Context) {
	var compclassJSON serializer.ComplexClassJSON
	if err := ctx.BindJSON(&compclassJSON); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	compclass, err := h.Repository.EditComplexClass(id, compclassJSON)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, serializer.CompClassToJSON(compclass))
}

func (h *Handler) AddToBigORequest(ctx *gin.Context) {
	bigorequest, created, err := h.Repository.GetBigORequestDraft(uint(h.Repository.GetUserID()))
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	bigorequest_id := bigorequest.ID

	complass_id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	err = h.Repository.AddToBigORequest(int(bigorequest_id), complass_id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else if errors.Is(err, repository.ErrAlreadyExists) {
			h.errorHandler(ctx, http.StatusConflict, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	status := http.StatusOK

	if created {
		ctx.Header("Location", fmt.Sprintf("/BigORequest/%v", bigorequest.ID))
		status = http.StatusCreated
	}

	creatorLogin, moderatorLogin, err := h.Repository.GetModeratorAndCreatorLogin(bigorequest)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(status, serializer.BigORequestToJSON(bigorequest, creatorLogin, moderatorLogin))
}

func (h *Handler) AddPhoto(ctx *gin.Context) {
	compclass_id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	file, err := ctx.FormFile("image")
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	compclass, err := h.Repository.AddPhoto(ctx, compclass_id, file)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "uploaded",
		"device": serializer.CompClassToJSON(compclass),
	})
}
