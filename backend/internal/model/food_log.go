package model

import "time"

type FoodLog struct {
	ID          int64     `json:"id"`
	TelegramID  int64     `json:"telegram_id"`
	Barcode     string    `json:"barcode"`
	ProductName string    `json:"product_name"`
	WeightGrams float64   `json:"weight_grams"`
	Kcal        float64   `json:"kcal"`
	Protein     float64   `json:"protein"`
	Fat         float64   `json:"fat"`
	Carbs       float64   `json:"carbs"`
	CreatedAt   time.Time `json:"created_at"`
}

type LogRequest struct {
	TelegramID  int64   `json:"telegram_id"`
	Barcode     string  `json:"barcode"`
	WeightGrams float64 `json:"weight_grams"`
}

type DailyStatus struct {
	Date            string  `json:"date"`
	DailyKcalTarget float64 `json:"daily_kcal_target"`
	ConsumedKcal    float64 `json:"consumed_kcal"`
	ConsumedProtein float64 `json:"consumed_protein"`
	ConsumedFat     float64 `json:"consumed_fat"`
	ConsumedCarbs   float64 `json:"consumed_carbs"`
	RemainingKcal   float64 `json:"remaining_kcal"`
}
