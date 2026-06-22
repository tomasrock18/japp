package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/tomasrock18/japp/backend/internal/model"
	"github.com/tomasrock18/japp/backend/internal/storage"
)

type LogHandler struct {
	logStorage     *storage.LogStorage
	productStorage *storage.MemoryStorage
}

func NewLogHandler(logStorage *storage.LogStorage, productStorage *storage.MemoryStorage) *LogHandler {
	return &LogHandler{
		logStorage:     logStorage,
		productStorage: productStorage,
	}
}

func (h *LogHandler) CreateLog(w http.ResponseWriter, r *http.Request) {
	var req model.LogRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "invalid request body"}`, http.StatusBadRequest)
		return
	}
	if req.TelegramID == 0 {
		http.Error(w, `{"error": "telegram id is required"}`, http.StatusBadRequest)
		return
	}
	if req.Barcode == "" {
		http.Error(w, `{"error": "barcode is required"}`, http.StatusBadRequest)
		return
	}
	if req.WeightGrams <= 0 {
		http.Error(w, `{"error": "weight_grams is required"}`, http.StatusBadRequest)
		return
	}

	product, err := h.productStorage.GetProduct(req.Barcode)
	if err != nil {
		http.Error(w, `{"error": "product not found"}`, http.StatusNotFound)
		return
	}

	multiplier := req.WeightGrams / 100.0
	foodLog := model.FoodLog{
		TelegramID:  req.TelegramID,
		Barcode:     req.Barcode,
		ProductName: product.Name,
		WeightGrams: req.WeightGrams,
		Kcal:        multiplier * product.KcalPer100g,
		Protein:     multiplier * product.ProteinPer100g,
		Fat:         multiplier * product.FatPer100g,
		Carbs:       multiplier * product.CarbsPer100g,
	}

	createdLog, err := h.logStorage.CreateLog(foodLog)
	if err != nil {
		http.Error(w, `{"error": "failed to create log"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(createdLog); err != nil {
		slog.Warn("Failed to encode log", "error", err)
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *LogHandler) GetDailyStats(w http.ResponseWriter, r *http.Request) {
	telegramIDStr := chi.URLParam(r, "telegram_id")
	telegramID, err := strconv.ParseInt(telegramIDStr, 10, 64)
	if err != nil {
		http.Error(w, `{"error": "invalid telegram_id"`, http.StatusBadRequest)
		return
	}

	dateStr := r.URL.Query().Get("date")

	var date time.Time
	if dateStr == "" {
		date = time.Now()
	} else {
		date, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			http.Error(w, `{"error": "invalid date format"}`, http.StatusBadRequest)
			return
		}
	}

	logs, err := h.logStorage.GetLogsByUserAndDate(telegramID, date)
	if err != nil {
		http.Error(w, `{"error": "failed to get logs"}`, http.StatusInternalServerError)
		return
	}

	var totalKcal, totalProtein, totalFat, totalCarbs float64
	for _, log := range logs {
		totalKcal += log.Kcal
		totalProtein += log.Protein
		totalFat += log.Fat
		totalCarbs += log.Carbs
	}

	dailyTarget := 2000.0
	stats := model.DailyStatus{
		Date:            date.Format("2006-01-02"),
		DailyKcalTarget: dailyTarget,
		ConsumedKcal:    totalKcal,
		ConsumedProtein: totalProtein,
		ConsumedFat:     totalFat,
		ConsumedCarbs:   totalCarbs,
		RemainingKcal:   dailyTarget - totalKcal,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(stats); err != nil {
		slog.Warn("Failed to encode daily stats", "error", err)
	}
}

func (h *LogHandler) GetUserLogs(w http.ResponseWriter, r *http.Request) {
	telegramIDStr := chi.URLParam(r, "telegram_id")

	telegramID, err := strconv.ParseInt(telegramIDStr, 10, 64)
	if err != nil {
		http.Error(w, `{"error": "Invalid telegram id"}`, http.StatusBadRequest)
		return
	}

	limit := 50
	limitStr := r.URL.Query().Get("limit")
	if limitStr != "" {
		limitValue, err := strconv.Atoi(limitStr)
		if err != nil {
			http.Error(w, `{"error": "Invalid limit parameter"}`, http.StatusBadRequest)
			return
		}
		if limitValue > 100 {
			http.Error(w, `{"error": "limit should be lesser than 100"}`, http.StatusBadRequest)
			return
		}
		limit = limitValue
	}

	offset := 0
	offsetStr := r.URL.Query().Get("offset")
	if limitStr != "" {
		offsetValue, err := strconv.Atoi(offsetStr)
		if err != nil {
			http.Error(w, `{"error": "Invalid offset prarmeter"}`, http.StatusBadRequest)
			return
		}
		if offsetValue < 0 {
			http.Error(w, `{"error": "offset should be greater than 0"}`, http.StatusBadRequest)
			return
		}
		offset = offsetValue
	}

	logs, err := h.logStorage.GetLogsByUser(telegramID, limit, offset)
	if err != nil {
		http.Error(w, `{"error": "Failed to extract user logs"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(logs); err != nil {
		slog.Warn("Failed to encode user logs", "error", err)
	}

}
