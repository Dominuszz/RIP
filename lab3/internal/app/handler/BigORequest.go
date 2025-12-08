package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"lab3/internal/app/ds"
	"lab3/internal/app/repository"
	"lab3/internal/app/serializer"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// GetAllBigORequests godoc
// @Summary Получить список заявок на расчёт
// @Description Возвращает заявки с возможностью фильтрации по датам и статусу
// @Tags bigorequests
// @Produce json
// @Param from-date query string false "Начальная дата (YYYY-MM-DD)"
// @Param to-date query string false "Конечная дата (YYYY-MM-DD)"
// @Param status query string false "Статус заявки"
// @Success 200 {array} serializer.BigORequestJSON "Список заявок"
// @Failure 400 {object} map[string]string "Неверный формат даты"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /bigorequest/all-bigo_requests [get]
func (h *Handler) GetAllBigORequests(ctx *gin.Context) {
	fromDate := ctx.Query("from-date")
	var from = time.Time{}
	var to = time.Time{}
	if fromDate != "" {
		from1, err := time.Parse("2006-01-02", fromDate)
		if err != nil {
			h.errorHandler(ctx, http.StatusBadRequest, err)
			return
		}
		from = from1
	}
	fmt.Println(fromDate)

	toDate := ctx.Query("to-date")
	if toDate != "" {
		to1, err := time.Parse("2006-01-02", toDate)
		if err != nil {
			h.errorHandler(ctx, http.StatusBadRequest, err)
			return
		}
		to = to1
	}

	status := ctx.Query("status")

	bigorequests, err := h.Repository.GetAllBigORequests(from, to, status)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	bigorequests = h.filterAuthorizedBigorequests(bigorequests, ctx)
	resp := make([]serializer.BigORequestJSON, 0, len(bigorequests))
	for _, c := range bigorequests {
		creatorLogin, moderatorLogin, err := h.Repository.GetModeratorAndCreatorLogin(c)
		if err != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
			return
		}
		resp = append(resp, serializer.BigORequestToJSON(c, creatorLogin, moderatorLogin))
	}
	ctx.JSON(http.StatusOK, resp)
}

// GetBigORequestCart godoc
// @Summary Получить корзину расчёта
// @Description Возвращает информацию о текущей заявке-черновике на расчёт пользователя
// @Tags bigorequests
// @Produce json
// @Success 200 {object} map[string]interface{} "Данные корзины заявки-черновика"
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /bigorequest/bigorequest-cart [get]
func (h *Handler) GetBigORequestCart(ctx *gin.Context) {
	userID, err := getUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"id":              -1,
			"compclass_count": 0,
		})
		return
	}
	compclass_count := h.Repository.GetBigORequestCount(userID)

	if compclass_count == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"status":          "no_draft",
			"compclass_count": compclass_count,
		})
		return
	}

	bigorequest, err := h.Repository.CheckCurrentBigORequestDraft(userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotAllowed) {
			h.errorHandler(ctx, http.StatusUnauthorized, err)
		} else if errors.Is(err, repository.ErrNoDraft) {
			ctx.JSON(http.StatusOK, gin.H{
				"status":          "no_draft",
				"compclass_count": 0,
			})
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":              bigorequest.ID,
		"compclass_count": compclass_count,
	})
}

// GetBigORequest godoc
// @Summary Получить заявку по ID
// @Description Возвращает полную информацию о заявке
// @Tags bigorequests
// @Produce json
// @Param id path int true "ID заявки"
// @Success 200 {object} map[string]interface{} "Данные заявки с классами сложности"
// @Failure 400 {object} map[string]string "Неверный ID"
// @Failure 403 {object} map[string]string "Доступ запрещен"
// @Failure 404 {object} map[string]string "Заявка не найдено"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /bigorequest/{id} [get]
func (h *Handler) GetBigORequest(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	compclasses, bigorequest, err := h.Repository.GetBigORequestComplexClasses(id)
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

	resp := make([]serializer.ComplexClassJSON, 0, len(compclasses))
	for _, r := range compclasses {
		resp = append(resp, serializer.CompClassToJSON(r))
	}

	creatorLogin, moderatorLogin, err := h.Repository.GetModeratorAndCreatorLogin(bigorequest)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"bigorequest": serializer.BigORequestToJSON(bigorequest, creatorLogin, moderatorLogin),
		"compclasses": resp,
	})
}

// FormBigORequest godoc
// @Summary Сформировать заявку
// @Description Переводит заявку в статус "formed"
// @Tags bigorequests
// @Produce json
// @Param id path int true "ID заявки"
// @Success 200 {object} serializer.BigORequestJSON "Сформированная заявка"
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 403 {object} map[string]string "Доступ запрещен"
// @Failure 404 {object} map[string]string "Заявка не найдена"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /bigorequest/{id}/form-bigorequest [put]
func (h *Handler) FormBigORequest(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	status := "сформирован"

	bigorequest, err := h.Repository.FormBigORequest(id, status)
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

// EditBigORequest godoc
// @Summary Изменить заявку
// @Description Обновляет данные заявки
// @Tags bigorequests
// @Accept json
// @Produce json
// @Param id path int true "ID заявки"
// @Param bigorequest body serializer.BigORequestJSON true "Новые данные заявки"
// @Success 200 {object} serializer.BigORequestJSON "Обновленная заявка"
// @Failure 400 {object} map[string]string "Неверные данные"
// @Failure 404 {object} map[string]string "Заявка не найдена"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /bigorequest/{id}/edit-bigorequest [put]
func (h *Handler) EditBigORequest(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	var bigorequestJSON serializer.BigORequestJSON
	if err := ctx.BindJSON(&bigorequestJSON); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	bigorequest, err := h.Repository.EditBigORequest(id, bigorequestJSON)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
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

// DeleteBigORequest godoc
// @Summary Удалить заявку
// @Description Выполняет логическое удаление заявки
// @Tags bigorequests
// @Produce json
// @Param id path int true "ID заявки"
// @Success 200 {object} map[string]string "Статус удаления"
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 403 {object} map[string]string "Доступ запрещен"
// @Failure 404 {object} map[string]string "Заявка не найдена"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /bigorequest/{id}/delete-bigorequest [delete]
func (h *Handler) DeleteBigORequest(ctx *gin.Context) {
	idStr := ctx.Param("id")
	bigorequest_id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	status := "удален"

	_, err = h.Repository.FormBigORequest(bigorequest_id, status)
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

	ctx.JSON(http.StatusOK, gin.H{"message": "BigO Request deleted"})
}

// FinishBigORequest godoc
// @Summary Завершить заявку
// @Description Изменяет статус заявки (только для модераторов) и запускает асинхронный расчет
// @Tags bigorequests
// @Accept json
// @Produce json
// @Param id path int true "ID заявки"
// @Param status body serializer.StatusJSON true "Новый статус"
// @Success 200 {object} serializer.BigORequestJSON "Результат модерации"
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 403 {object} map[string]string "Доступ запрещен"
// @Failure 404 {object} map[string]string "Заявка не найдена"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /bigorequest/{id}/finish-bigorequest [put]
func (h *Handler) FinishBigORequest(ctx *gin.Context) {
	userID, err := getUserID(ctx)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	var statusJSON serializer.StatusJSON
	if err := ctx.BindJSON(&statusJSON); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	user, err := h.Repository.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	if !user.IsModerator {
		h.errorHandler(ctx, http.StatusForbidden, errors.New("требуются права модератора"))
		return
	}

	// Завершаем заявку (меняем статус на "выполнен" или "отклонен")
	bigorequest, err := h.Repository.FinishBigORequest(id, statusJSON.Status, userID)
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

	// Если статус "выполнен", запускаем асинхронный расчет
	if statusJSON.Status == "выполнен" {
		go h.startAsyncCalculation(id)
	}

	creatorLogin, moderatorLogin, err := h.Repository.GetModeratorAndCreatorLogin(bigorequest)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, serializer.BigORequestToJSON(bigorequest, creatorLogin, moderatorLogin))
}

// startAsyncCalculation запускает асинхронный расчет через Django сервис
func (h *Handler) startAsyncCalculation(requestID int) {
	logrus.Infof("Запуск асинхронного расчета для заявки %d", requestID)

	// Получение данных заявки для расчета
	compclasses, _, err := h.Repository.GetBigORequestComplexClasses(requestID)
	if err != nil {
		logrus.Errorf("Ошибка получения данных заявки %d: %v", requestID, err)
		return
	}

	// Получение размеров массивов для каждого класса
	compclassRequests, err := h.Repository.GetComplexClassesBigORequests(requestID)
	if err != nil {
		logrus.Errorf("Ошибка получения размеров массивов для заявки %d: %v", requestID, err)
		return
	}

	// Создание map для быстрого поиска размеров массивов
	arraySizes := make(map[uint]uint)
	for _, cr := range compclassRequests {
		arraySizes[cr.ComplexClassID] = cr.ArraySize
	}

	// Подготовка данных для асинхронного сервиса
	compclassesData := make([]map[string]interface{}, 0, len(compclasses))
	for _, cc := range compclasses {
		arraySize := arraySizes[uint(cc.ID)]
		if arraySize == 0 {
			arraySize = 1000 // Значение по умолчанию
		}

		compclassesData = append(compclassesData, map[string]interface{}{
			"complexity": cc.Complexity,
			"degree":     cc.Degree,
			"array_size": arraySize,
		})
	}

	// Формирование запроса к асинхронному сервису Django
	calcRequest := map[string]interface{}{
		"request_id":  requestID,
		"compclasses": compclassesData,
	}

	// URL Django сервиса
	asyncServiceURL := "http://localhost:8000/api/calculate/"

	jsonData, err := json.Marshal(calcRequest)
	if err != nil {
		logrus.Errorf("Ошибка сериализации запроса для заявки %d: %v", requestID, err)
		return
	}

	// Отправка запроса в асинхронный сервис
	resp, err := http.Post(asyncServiceURL, "application/json",
		strings.NewReader(string(jsonData)))
	if err != nil {
		logrus.Errorf("Ошибка вызова асинхронного сервиса для заявки %d: %v", requestID, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		body, _ := io.ReadAll(resp.Body)
		logrus.Errorf("Асинхронный сервис вернул ошибку для заявки %d: %d - %s",
			requestID, resp.StatusCode, string(body))
		return
	}

	logrus.Infof("Асинхронный расчет для заявки %d успешно запущен", requestID)
}

// UpdateCalculation godoc
// @Summary Обновить результаты расчета (для асинхронного сервиса)
// @Description Принимает результаты расчета от асинхронного сервиса
// @Tags bigorequests
// @Accept json
// @Produce json
// @Param id path int true "ID заявки"
// @Param secret-key header string true "Секретный ключ для авторизации"
// @Param result body map[string]interface{} true "Результаты расчета"
// @Success 200 {object} serializer.BigORequestJSON "Заявка обновлена"
// @Failure 400 {object} map[string]string "Неверный запрос или ключ"
// @Failure 403 {object} map[string]string "Доступ запрещен"
// @Failure 404 {object} map[string]string "Заявка не найдена"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /bigorequest/{id}/update-calculation [put]
// UpdateCalculation godoc
// @Summary Обновить результаты расчета (для асинхронного сервиса)
// @Description Принимает результаты расчета от асинхронного сервиса
// @Tags bigorequests
// @Accept json
// @Produce json
// @Param id path int true "ID заявки"
// @Param result body map[string]interface{} true "Результаты расчета"
// @Success 200 {object} serializer.BigORequestJSON "Заявка обновлена"
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 404 {object} map[string]string "Заявка не найдена"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /bigorequest/{id}/update-calculation [put]
func (h *Handler) UpdateCalculation(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	// Проверка существования заявки
	bigorequest, err := h.Repository.GetSingleBigORequest(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	// Проверка, что заявка находится в статусе "выполнен"
	if bigorequest.Status != "выполнен" {
		h.errorHandler(ctx, http.StatusBadRequest,
			errors.New("расчет возможен только для заявок со статусом 'выполнен'"))
		return
	}

	// Чтение данных расчета
	var calcResult map[string]interface{}
	if err := ctx.BindJSON(&calcResult); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	// Извлечение данных из результата
	calculatedTime, ok1 := calcResult["calculated_time"].(float64)
	calculatedComplexity, ok2 := calcResult["calculated_complexity"].(string)
	success, ok3 := calcResult["success"].(bool)

	if !ok1 || !ok2 || !ok3 {
		h.errorHandler(ctx, http.StatusBadRequest,
			errors.New("неверный формат результатов расчета"))
		return
	}

	// Обновление заявки с результатами расчета
	bigorequest.CalculatedTime = calculatedTime
	bigorequest.CalculatedComplexity = calculatedComplexity

	// Если расчет неуспешен, меняем статус на "отклонен"
	if !success {
		bigorequest.Status = "отклонен"
	}

	err = h.Repository.UpdateBigORequestResult(bigorequest)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	// Получение обновленной заявки
	updatedRequest, err := h.Repository.GetSingleBigORequest(id)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	// Получение логинов для ответа
	creatorLogin, moderatorLogin, err := h.Repository.GetModeratorAndCreatorLogin(updatedRequest)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, serializer.BigORequestToJSON(updatedRequest, creatorLogin, moderatorLogin))
}

func (h *Handler) filterAuthorizedBigorequests(bigorequests []ds.BigORequest, ctx *gin.Context) []ds.BigORequest {
	userID, err := getUserID(ctx)
	if err != nil {
		return []ds.BigORequest{}
	}

	user, err := h.Repository.GetUserByID(userID)
	if err == repository.ErrNotFound {
		return []ds.BigORequest{}
	}
	if err != nil {
		return []ds.BigORequest{}
	}

	if user.IsModerator {
		return bigorequests
	}

	var userBigorequests []ds.BigORequest
	for _, bigorequest := range bigorequests {
		fmt.Println(bigorequest.ID)
		if bigorequest.CreatorID == userID {
			userBigorequests = append(userBigorequests, bigorequest)
		}
	}

	return userBigorequests

}

func (h *Handler) hasAccessToBigORequest(creatorID uuid.UUID, ctx *gin.Context) bool {
	userID, err := getUserID(ctx)
	if err != nil {
		return false
	}

	user, err := h.Repository.GetUserByID(userID)
	if err == repository.ErrNotFound {
		return false
	}
	if err != nil {
		return false
	}

	return creatorID == userID || user.IsModerator
}
