package handler

import (
	"errors"
	"lab3/internal/app/repository"
	"lab3/internal/app/serializer"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) DeleteCompClassFromBigORequest(ctx *gin.Context) {
	bigorequest_id, err := strconv.Atoi(ctx.Param("bigo_request_id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	compclass_id, err := strconv.Atoi(ctx.Param("compclass_id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	bigorequest, err := h.Repository.DeleteCompClassFromBigORequest(bigorequest_id, compclass_id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else if errors.Is(err, repository.ErrNotAllowed) {
			h.errorHandler(ctx, http.StatusForbidden, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	creatorLogin, moderatorLogin, err := h.Repository.GetModeratorAndCreatorLogin(bigorequest)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, serializer.BigORequestToJSON(bigorequest, creatorLogin, moderatorLogin))
}

func (h *Handler) EditCompClassFromBigORequest(ctx *gin.Context) {
	bigorequest_id, err := strconv.Atoi(ctx.Param("bigo_request_id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	compclass_id, err := strconv.Atoi(ctx.Param("compclass_id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	var CompClassRequestJSON serializer.CompClassRequestJSON
	if err := ctx.BindJSON(&CompClassRequestJSON); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	compclassrequest, err := h.Repository.EditCompClassFromBigORequest(bigorequest_id, compclass_id, CompClassRequestJSON)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, serializer.CompClassRequestToJSON(compclassrequest))
}
