package model

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type Product struct {
	Barcode        string    `json:"barcode"`
	Name           string    `json:"name"`
	KcalPer100g    float64   `json:"kcal_per_100g"`
	ProteinPer100g float64   `json:"protein_per_100g"`
	FatPer100g     float64   `json:"fat_per_100g"`
	CarbsPer100g   float64   `json:"carbs_per_100g"`
	CreatedBy      int64     `json:"created_by"`
	CreatedAt      time.Time `json:"created_at"`
}

func (product *Product) IsValid() error {
	if strings.TrimSpace(product.Barcode) == "" {
		return errors.New("barcode is empty")
	}

	if strings.TrimSpace(product.Name) == "" {
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
		if parsedValue, ok := value.(string); ok {
			product.Barcode = parsedValue
		} else {
			return fmt.Errorf("bad barcode value: %v", value)
		}
	case "name":
		if parsedValue, ok := value.(string); ok {
			product.Name = parsedValue
		} else {
			return fmt.Errorf("bad name value: %v", value)
		}
	case "kcal_per_100g":
		if parsedValue, ok := value.(float64); ok {
			product.KcalPer100g = parsedValue
		} else {
			return fmt.Errorf("bad kcal_per_100g value: %v", value)
		}
	case "protein_per_100g":
		if parsedValue, ok := value.(float64); ok {
			product.ProteinPer100g = parsedValue
		} else {
			return fmt.Errorf("bad protein_per_100g value: %v", value)
		}
	case "fat_per_100g":
		if parsedValue, ok := value.(float64); ok {
			product.FatPer100g = parsedValue
		} else {
			return fmt.Errorf("bad fat_per_100g value: %v", value)
		}
	case "carbs_per_100g":
		if parsedValue, ok := value.(float64); ok {
			product.CarbsPer100g = parsedValue
		} else {
			return fmt.Errorf("bad carbs_per_100g value: %v", value)
		}
	default:
		return fmt.Errorf("unknown field %s", field)
	}
	return nil
}
