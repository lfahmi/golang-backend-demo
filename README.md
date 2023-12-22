# golang-backend-demo
My demonstration for Golang Backend Developer. This project and documentation is done whitin 1 day.


How to Setup in Local Environment:
Setup Steps:
Install Go:

Download and install Golang from the official site.
Make sure the system PATH includes the bin directory of the Golang installation.
Install Dependencies:

Use the go get command to download and install the necessary libraries, such as Gin, Gorm

go get -u github.com/gin-gonic/gin
go get -u gorm.io/gorm
go get -u github.com/dgrijalva/jwt-go

Database Configuration:
Set the database connection in the config.go file according to your local settings or more easily by importing the sql schema in /doc/
including the user and privileges

Running the Application:
Use the go run main.go command from the terminal to run the application.
The application will run on localhost:8080 or the address specified in the configuration.

Data Migration and Embedding:
After the application is run, the creation of the database schema and dummy data will automatically be carried out.

By following these steps, you should be able to run and test the project in your local environment. Make sure that system requirements such as Golang are met.