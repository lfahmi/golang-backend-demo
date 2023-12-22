package main

import "github.com/jinzhu/gorm"

func seedRatingTable(db *gorm.DB) error {
	var models []Rating
	var count int
	if err := db.Unscoped().Model(&models).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	seeds := []*Rating{
		{
			UserID: 1,
			GiftID: 1,
			Rate:   5,
		},
		{
			UserID: 1,
			GiftID: 1,
			Rate:   1,
		},
		{
			UserID: 2,
			GiftID: 1,
			Rate:   4,
		},
		{
			UserID: 1,
			GiftID: 2,
			Rate:   4,
		},
		{
			UserID: 2,
			GiftID: 2,
			Rate:   2,
		},
		{
			UserID: 1,
			GiftID: 3,
			Rate:   5,
		},
		{
			UserID: 2,
			GiftID: 3,
			Rate:   1,
		},
	}

	for _, v := range seeds {
		if err := db.Create(v).Error; err != nil {
			return err
		}
	}
	return nil
}
