package main

import "github.com/jinzhu/gorm"

func seedUserTable(db *gorm.DB) error {
	var models []User
	var count int
	if err := db.Unscoped().Model(&models).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	seeds := []*User{
		{
			Model:    gorm.Model{ID: 1},
			Username: "user",
			Password: "user",
		},
		{
			Model:    gorm.Model{ID: 2},
			Username: "john",
			Password: "doe",
		},
	}

	for _, v := range seeds {
		if err := db.Create(v).Error; err != nil {
			return err
		}
	}
	return nil
}
