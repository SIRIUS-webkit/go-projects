package main

import (
	"net/http"

	"github.com/SIRIUS-webkit/crud-app/internal/database"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)	

type RegisterUserRequest struct {
	Email string `json:"email" binding:"required,email"`
	Name string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}

func (app *application) registerUser(c *gin.Context){
	var request RegisterUserRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}	
	
	request.Password = string(hashedPassword)

	user := database.User{
		Email: request.Email,
		Name: request.Name,
		Password: request.Password,
	}

	err = app.models.Users.Insert(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert user"})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
	
}


