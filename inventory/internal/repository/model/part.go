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

// Part представляет сущность детали в repository слое
type Part struct {
	UUID         string        `bson:"uuid"`
	Name         string        `bson:"name"`
	Description  string        `bson:"description"`
	Price        float64       `bson:"price"`
	Category     Category      `bson:"category"`
	Dimensions   *Dimensions   `bson:"dimensions,omitempty"`
	Manufacturer *Manufacturer `bson:"manufacturer,omitempty"`
	Tags         []string      `bson:"tags,omitempty"`
	CreatedAt    time.Time     `bson:"createdAt"`
	UpdatedAt    time.Time     `bson:"updatedAt"`
}

// Dimensions представляет размеры детали
type Dimensions struct {
	Length float64 `bson:"length"`
	Width  float64 `bson:"width"`
	Height float64 `bson:"height"`
	Weight float64 `bson:"weight"`
}

// Manufacturer представляет производителя детали
type Manufacturer struct {
	Name    string `bson:"name"`
	Country string `bson:"country"`
	Website string `bson:"website,omitempty"`
}
