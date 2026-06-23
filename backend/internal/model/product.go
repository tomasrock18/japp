package model

import (
	"errors"
	"fmt"
)

type Product struct {
	Barcode        string  `json:"barcode"`
	Name           string  `json:"name"`
	KcalPer100g    float64 `json:"kcal_per_100g"`
	ProteinPer100g float64 `json:"protein_per_100g"`
	FatPer100g     float64 `json:"fat_per_100g"`
	CarbsPer100g   float64 `json:"carbs_per_100g"`
	CreatedBy      string  `json:"created_by,omitempty"`
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

func (product *Product) UpdateField(field string, value any) error {
	switch field {
	case "barcode":
		product.Barcode = value.(string)
	case "name":
		product.Name = value.(string)
	case "kcal_per_100g":
		product.KcalPer100g = value.(float64)
	case "protein_per_100g":
		product.ProteinPer100g = value.(float64)
	case "fat_per_100g":
		product.FatPer100g = value.(float64)
	case "carbs_per_100g":
		product.CarbsPer100g = value.(float64)
	default:
		return fmt.Errorf("unknown field %s", field)
	}
	return nil
}
