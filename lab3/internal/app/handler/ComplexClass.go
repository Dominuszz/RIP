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

// GetComplexClasses godoc
// @Summary Получить список классов сложности
// @Description Возвращает все классы сложности или фильтрует по степени
// @Tags CompClasses
// @Produce json
// @Param search-degree query string false "Степень класса сложности для поиска"
// @Success 200 {array} serializer.ComplexClassJSON "Список классов сложности"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /complexclass [get]
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

// GetComplexClass godoc
// @Summary Получить класс сложности по ID
// @Description Возвращает информацию о классе сложности по его идентификатору
// @Tags CompClasses
// @Produce json
// @Param id path int true "ID класса сложности"
// @Success 200 {object} serializer.ComplexClassJSON "Данные класса сложности"
// @Failure 400 {object} map[string]string "Неверный ID"
// @Failure 404 {object} map[string]string "Устройство не найдено"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /complexclass/{id} [get]
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

// CreateComplexClass godoc
// @Summary Создать новый класс сложности
// @Description Создает новый класс сложности и возвращает его данные
// @Tags CompClasses
// @Accept json
// @Produce json
// @Param device body serializer.ComplexClassJSON true "Данные нового класса сложности"
// @Success 201 {object} serializer.ComplexClassJSON "Созданный класс сложности"
// @Failure 400 {object} map[string]string "Неверные данные"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /complexclass/create-compclass [post]
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

// DeleteComplexClas godoc
// @Summary Удалить класс сложности
// @Description Выполняет логическое удаление класса сложности по ID
// @Tags CompClasses
// @Produce json
// @Param id path int true "ID класса сложности"
// @Success 200 {object} map[string]string "Статус удаления"
// @Failure 400 {object} map[string]string "Неверный ID"
// @Failure 404 {object} map[string]string "Класс сложности не найден"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /complexclass/{id}/delete-compclass [delete]
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

// EditComplexClass godoc
// @Summary Изменить данные класса сложности
// @Description Обновляет информацию о классе сложности по ID
// @Tags CompClasses
// @Accept json
// @Produce json
// @Param id path int true "ID класса сложности"
// @Param compclass body serializer.ComplexClassJSON  true "Новые данные класса сложности"
// @Success 200 {object} serializer.ComplexClassJSON  "Обновленный класс сложности"
// @Failure 400 {object} map[string]string "Неверные данные"
// @Failure 404 {object} map[string]string "Класс сложности не найден"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /complexclass/{id}/edit-compclass [put]
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

// AddToBigORequest godoc
// @Summary Добавить класс сложности в расчёт
// @Description Добавляет класс сложности в заявку-черновик пользователя
// @Tags CompClasses
// @Produce json
// @Param id path int true "ID класса сложности"
// @Success 200 {object} serializer.BigORequestJSON "Расчёт с добавленным классом сложности"
// @Success 201 {object} serializer.BigORequestJSON "Создан новый расчёт"
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 404 {object} map[string]string "Класс сложности не найден"
// @Failure 409 {object} map[string]string "Класс сложности уже в расчёте"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /complexclass/{id}/add-to-bigorequest [post]
func (h *Handler) AddToBigORequest(ctx *gin.Context) {
	userID, err := getUserID(ctx)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	bigorequest, created, err := h.Repository.GetBigORequestDraft(userID)
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

// AddPhoto godoc
// @Summary Загрузить изображение устройства
// @Description Загружает изображение для класса сложности и возвращает обновленные данные
// @Tags CompClasses
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "ID класса сложности"
// @Param image formData file true "Изображение класса сложности"
// @Success 200 {object} map[string]interface{} "Статус загрузки и данные класса сложности"
// @Failure 400 {object} map[string]string "Неверный запрос или файл"
// @Failure 404 {object} map[string]string "Класс сложности не найден"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /complexclass/{id}/add-photo [post]
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
		"status":    "uploaded",
		"compclass": serializer.CompClassToJSON(compclass),
	})
}
