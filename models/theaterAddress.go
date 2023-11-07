package models

import "gorm.io/gorm"

type TheaterAddress struct {
	gorm.Model
	TheaterId int64 `json:"theaterId"`
	AddressId int64 `json:"addressId"`
}
