package main

import (
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// JWT middleware
func authMiddleware(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" || strings.Index(tokenString, "Bearer ") != 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}

	tokenString = tokenString[len("Bearer "):]

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwt_secret, nil
	})

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		c.Set("claims", claims)
		c.Next()
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
	}
}

func getCurrentUserId(c *gin.Context) (uint, bool) {
	// Get the user ID from the JWT claims
	claimsObj, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "JWT claims not found in context"})
		return 0, false
	}
	claims, ok := claimsObj.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "JWT claims not found in context"})
		return 0, false
	}
	userID, ok := claims["uid"].(float64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to extract user ID from JWT claims"})
		return 0, false
	}
	return uint(userID), true
}
