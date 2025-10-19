package model

import (
	"time"
)

// Category представляет категорию детали
type Category int32

const (
	CategoryUnspecified Category = iota
	CategoryEngine
	CategoryFuel
	CategoryPorthole
	CategoryWing
)

// Part представляет сущность детали в сервисном слое
type Part struct {
	UUID          string
	Name          string
	Description   string
	Price         float64
	StockQuantity int32
	Category      Category
	Dimensions    *Dimensions
	Manufacturer  *Manufacturer
	Tags          []string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// Dimensions представляет размеры детали
type Dimensions struct {
	Length float64
	Width  float64
	Height float64
	Weight float64
}

// Manufacturer представляет производителя детали
type Manufacturer struct {
	Name    string
	Country string
	Website string
}

// Filter представляет фильтр для поиска деталей
type Filter struct {
	UUIDs                 []string
	Names                 []string
	Categories            []Category
	ManufacturerCountries []string
	Tags                  []string
}
