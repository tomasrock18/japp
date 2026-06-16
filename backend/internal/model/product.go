package model

import "errors"

type Product struct {
	Barcode        string  `json:"barcode"`
	Name           string  `json:"name"`
	KcalPer100g    float64 `json:"kcal_per_100g"`
	ProteinPer100g float64 `json:"protein_per_100g"`
	FatPer100g     float64 `json:"fat_per_100g"`
	CarbsPer100g   float64 `json:"carbs_per_100g"`
}

func (product *Product) IsValid() error {
	if product.Barcode == "" {
		return errors.New("barcode is empty")
	}

	if product.Name == "" {
		return errors.New("name is empty")
	}

	if product.KcalPer100g < 0 {
		return errors.New("kcal_per_100g should be greater than zero")
	}

	if product.ProteinPer100g < 0 {
		return errors.New("protein_per_100g should be greater than zero")
	}

	if product.FatPer100g < 0 {
		return errors.New("fat_per_100g should be greater than zero")
	}

	if product.CarbsPer100g < 0 {
		return errors.New("carbs_per_100g should be greater than zero")
	}

	return nil
}
