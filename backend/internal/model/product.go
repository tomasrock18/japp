package model

type Product struct {
	Barcode        string  `json:"barcode"`
	Name           string  `json:"name"`
	KcalPer100g    float64 `json:"kcal_per_100g"`
	ProteinPer100g float64 `json:"protein_per_100g"`
	FatPer100g     float64 `json:"fat_per_100g"`
	CarbsPer100g   float64 `json:"carbs_per_100g"`
}
