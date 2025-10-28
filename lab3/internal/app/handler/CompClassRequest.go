package handler

import (
	"errors"
	"lab3/internal/app/repository"
	"lab3/internal/app/serializer"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// DeleteCompClassFromBigORequest godoc
// @Summary Удалить класс сложности из заявки
// @Description Удаляет связь класса сложности и заявки
// @Tags CompClassRequest
// @Produce json
// @Param compclass_id path int true "ID класса сложности"
// @Param bigo_request_id path int true "ID заявки"
// @Success 200 {object} serializer.BigORequestJSON "Обновленная заявка"
// @Failure 400 {object} map[string]string "Неверные ID"
// @Failure 403 {object} map[string]string "Доступ запрещен"
// @Failure 404 {object} map[string]string "Не найдено"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /compclassrequest/{compclass_id}/{bigo_request_id} [delete]
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

// EditCompClassFromBigORequest godoc
// @Summary Изменить данные класса сложности в заявке
// @Description Обновляет параметры класса сложности в конкретной заявке
// @Tags CompClassRequest
// @Accept json
// @Produce json
// @Param compclass_id path int true "ID класса сложности"
// @Param bigo_request_id path int true "ID заявки"
// @Param data body serializer.CompClassRequestJSON true "Новые данные"
// @Success 200 {object} serializer.CompClassRequestJSON "Обновленные данные"
// @Failure 400 {object} map[string]string "Неверные данные"
// @Failure 404 {object} map[string]string "Не найдено"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /compclassrequest/{compclass_id}/{bigo_request_id} [put]
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
