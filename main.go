package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	db *gorm.DB
)

func init() {
	var err error
	db, err = gorm.Open("mysql", ""+DB_USER+":"+DB_PASS+"@"+DB_HOST+"/"+DB_NAME+"?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Fatal(err)
	}

	// Migration and Seeders
	db.AutoMigrate(&Gift{}, &User{}, &Rating{})
	seedUserTable(db)
	seedGiftTable(db)
	seedRatingTable(db)
}

func main() {
	r := gin.Default()

	// Initialize routes
	initializeRoutes(r)

	// Start the server
	r.Run(fmt.Sprintf(":%d", listen_port))
}

// Initialize routes
func initializeRoutes(r *gin.Engine) {
	// Public route for user login
	r.POST("/login", loginHandler)

	// Protected routes with JWT authentication middleware
	authenticated := r.Group("/")
	authenticated.Use(authMiddleware)

	// Gifts endpoints
	r.GET("/gifts", getGiftsHandler)
	r.GET("/gifts/:id", getGiftDetailsHandler)
	authenticated.POST("/gifts", createGiftHandler)
	authenticated.PUT("/gifts/:id", updateGiftHandler)
	authenticated.PATCH("/gifts/:id", patchGiftHandler)
	authenticated.DELETE(("/gifts/:id"), deleteGiftHandler)
	authenticated.POST("/gifts/:id/redeem", redeemGiftHandler)
	authenticated.POST("/gifts/:id/rating", rateGiftHandler)
}
