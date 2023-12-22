package main

import (
	"encoding/json"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// Handler to create a new gift
func createGiftHandler(c *gin.Context) {
	userID, ok := getCurrentUserId(c)
	if !ok {
		return
	}

	var newGift Gift
	if err := c.ShouldBindJSON(&newGift); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	newGift.CreatedBy = userID

	// Save the new gift to the database
	if err := db.Create(&newGift).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create gift"})
		return
	}

	c.JSON(http.StatusCreated, newGift)
}

// Handler to update an existing gift
func updateGiftHandler(c *gin.Context) {
	giftID := c.Param("id")

	// Get the user ID from the JWT claims
	userID, ok := getCurrentUserId(c)
	if !ok {
		return
	}

	var updatedGift Gift

	if err := c.ShouldBindJSON(&updatedGift); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Ensure that all attributes are filled
	if updatedGift.Name == "" || updatedGift.Description == "" || updatedGift.Stock < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "All attributes are required"})
		return
	}

	var existingGift Gift
	if err := db.Where("id = ? AND created_by = ?", giftID, userID).Find(&existingGift).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Gift not found"})
		return
	}

	// using map to force update even if it set stock to 0
	var updateData = map[string]interface{}{
		"ID":          giftID,
		"Name":        updatedGift.Name,
		"Description": updatedGift.Description,
		"Stock":       updatedGift.Stock,
		"Cost":        updatedGift.Cost,
	}

	// Update the existing gift
	if err := db.Model(&existingGift).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update gift"})
		return
	}

	c.JSON(http.StatusOK, existingGift)
}

// Handler to partially update an existing gift
func patchGiftHandler(c *gin.Context) {
	giftID := c.Param("id")

	// Get the user ID from the JWT claims
	userID, ok := getCurrentUserId(c)
	if !ok {
		return
	}

	var existingGift Gift
	if err := db.Where("id = ? AND created_by = ?", giftID, userID).Find(&existingGift).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Gift not found"})
		return
	}

	bodyStr, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// make map of updated attributes
	updatedAttrs := make(map[string]interface{})
	json.Unmarshal(bodyStr, &updatedAttrs)

	// ensure attributes name is valid
	for k := range updatedAttrs {
		if k != "name" && k != "description" && k != "stock" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}
	}

	// Perform the partial update
	if err := db.Model(&existingGift).Updates(updatedAttrs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update gift"})
		return
	}

	c.JSON(http.StatusOK, existingGift)
}

// Handler to delete an existing gift
func deleteGiftHandler(c *gin.Context) {
	giftID := c.Param("id")

	// Get the user ID from the JWT claims
	userID, ok := getCurrentUserId(c)
	if !ok {
		return
	}

	var existingGift Gift
	if err := db.Where("id = ? AND created_by = ?", giftID, userID).Find(&existingGift).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Gift not found"})
		return
	}

	// Delete the existing gift
	if err := db.Delete(&existingGift).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete gift"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Gift deleted successfully"})
}

// Handler to redeem gift
func redeemGiftHandler(c *gin.Context) {
	giftID := c.Param("id")
	var redeem Redeem

	if err := c.ShouldBindJSON(&redeem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Ensure that all attributes are filled
	if redeem.Quantity <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Quantity is required and must not be negative"})
		return
	}

	var existingGift Gift
	if err := db.First(&existingGift, giftID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Gift not found"})
		return
	}

	var updatedStock = existingGift.Stock - redeem.Quantity

	if updatedStock < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not enough stock of the item"})
		return
	}

	// Perform the partial update
	if err := db.Model(&existingGift).Update("Stock", updatedStock).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update gift"})
		return
	}

	c.JSON(http.StatusOK, existingGift)
}

func rateGiftHandler(c *gin.Context) {
	giftIDStr := c.Param("id")
	giftId, _ := strconv.ParseUint(giftIDStr, 10, 64)

	// Get the user ID from the JWT claims
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "JWT claims not found in context"})
		return
	}
	userID, ok := claims.(jwt.MapClaims)["uid"].(float64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to extract user ID from JWT claims"})
		return
	}

	bodyStr, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// make map of updated attributes
	updatedAttrs := make(map[string]interface{})
	json.Unmarshal(bodyStr, &updatedAttrs)

	// ensure attributes is valid
	for k := range updatedAttrs {
		if k != "rate" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}
	}
	rate, ok := updatedAttrs["rate"].(float64)
	if !ok || rate < 0 || rate > 5 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// ensure gift exists
	var existingGift Gift
	if err := db.First(&existingGift, giftId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Gift not found"})
		return
	}

	var rating Rating
	rating.UserID = uint(userID)
	rating.GiftID = uint(giftId)
	rating.Rate = int(rate)

	// Save the rating to the database
	if err := db.Create(&rating).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to rate the product"})
		return
	}

	c.JSON(http.StatusCreated, rating)
}

func getGiftDetailsHandler(c *gin.Context) {
	giftID := c.Param("id")

	// Fetch the gift details from the database, including average rating
	var gift struct {
		Gift
		RatingAverage float64 `json:"rating_average"`
		RatingsCount  int     `json:"ratings_count"`
	}

	if err := db.Table("gifts").
		Select("gifts.*, COALESCE(AVG(ratings.rate), 0) as rating_average, COUNT(ratings.id) as ratings_count").
		Joins("LEFT JOIN ratings ON gifts.id = ratings.gift_id").
		Where("gifts.id = ? AND gifts.deleted_at IS NULL", giftID). // Exclude soft-deleted records
		Group("gifts.id").
		First(&gift).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Gift not found"})
		return
	}

	gift.RatingAverage = math.Round(gift.RatingAverage*2) / 2

	c.JSON(http.StatusOK, gift)
}

func getGiftsHandler(c *gin.Context) {
	// Extract pagination parameters from query parameters
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page <= 0 {
		page = 1
	}

	perPage, err := strconv.Atoi(c.Query("per_page"))
	if err != nil || perPage <= 0 {
		perPage = 10 // Default per page count
	}

	offset := (page - 1) * perPage

	var gifts []struct {
		Gift
		RatingAverage float64 `json:"rating_average"`
		RatingsCount  int     `json:"ratings_count"`
	}

	// Fetch gifts with average rating from the database

	orderBy := c.Query("order_by")
	// constraint order
	order := strings.ToUpper(c.Query("order"))
	if order != "ASC" && order != "DESC" {
		order = ""
	}

	// build the query
	qry := db.Table("gifts").
		Select("gifts.*, COALESCE(AVG(ratings.rate), 0) as rating_average, COUNT(ratings.id) as ratings_count").
		Joins("LEFT JOIN ratings ON gifts.id = ratings.gift_id").
		Where("gifts.deleted_at IS NULL"). // Exclude soft-deleted records
		Group("gifts.id")
	if order != "" && orderBy == "rating" {
		qry = qry.Order("rating_average " + order)
	} else if order != "" && orderBy == "newest" {
		qry = qry.Order("created_at " + order)
	}
	qry = qry.Offset(offset).
		Limit(perPage)

	// fetch gifts using the query
	if err := qry.Find(&gifts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch gifts with average rating"})
		return
	}

	// make RatingAverage value become nearest .5
	for i, gift := range gifts {
		nearestRating := math.Round(gift.RatingAverage*2) / 2
		gifts[i].RatingAverage = nearestRating
	}

	c.JSON(http.StatusOK, gifts)
}
