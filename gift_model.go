package main

import "github.com/jinzhu/gorm"

// Gift model
type Gift struct {
	gorm.Model
	Name        string `json:"name"`
	Description string `json:"description"`
	Cost        int    `json:"cost"`
	Stock       int    `json:"stock"`
	CreatedBy   uint   `json:"created_by"`
}
