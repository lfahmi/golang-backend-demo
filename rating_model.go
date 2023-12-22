package main

import "github.com/jinzhu/gorm"

// Rating model
type Rating struct {
	gorm.Model
	UserID uint `json:"user_id"`
	GiftID uint `json:"gift_id"`
	Rate   int  `json:"rate"`
}
