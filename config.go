package main

var (
	listen_port = 8080
	jwt_secret  = []byte("oasjdoauwnaisudhaowndoaunsdouanw") // Change this with a secure random key in production and store it securely
	DB_USER     = "rgbtest"
	DB_PASS     = "rgbtest"
	DB_HOST     = "tcp(localhost:3306)"
	DB_NAME     = "fahmi-rgb-golang-test"
)
