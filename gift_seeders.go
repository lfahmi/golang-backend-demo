package main

import "github.com/jinzhu/gorm"

func seedGiftTable(db *gorm.DB) error {
	var models []Gift
	var count int
	if err := db.Unscoped().Model(&models).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	seeds := []*Gift{
		{
			Model:       gorm.Model{ID: 1},
			Name:        "Samsung galaxy 5g 4/64gb",
			Description: "Ukuran Layar 6.2 inci, Dual Edge Super AMOLED",
			Cost:        200000,
			Stock:       10,
			CreatedBy:   1,
		},
		{
			Model:       gorm.Model{ID: 2},
			Name:        "Xiaomi Redmi 5g 4/64gb",
			Description: "Ukuran Layar 6.4 inci, Dual Edge Super AMOLED",
			Cost:        180000,
			Stock:       8,
			CreatedBy:   1,
		},
		{
			Model:       gorm.Model{ID: 3},
			Name:        "Samsung galaxy 5g 6/128gb",
			Description: "Ukuran Layar 6.2 inci, Dual Edge Super AMOLED",
			Cost:        250000,
			Stock:       10,
			CreatedBy:   1,
		},
		{
			Model:       gorm.Model{ID: 4},
			Name:        "Xiaomi Redmi 5g 6/128gb",
			Description: "Ukuran Layar 6.4 inci, Dual Edge Super AMOLED",
			Cost:        230000,
			Stock:       8,
			CreatedBy:   1,
		},
		{
			Model:       gorm.Model{ID: 4},
			Name:        "Ipad 2020 WiFi",
			Description: "Ukuran Layar Besar dan performa tinggi",
			Cost:        480,
			Stock:       7,
			CreatedBy:   1,
		},
	}

	for _, v := range seeds {
		if err := db.Create(v).Error; err != nil {
			return err
		}
	}
	return nil
}
