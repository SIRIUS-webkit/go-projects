package main

import (
	"net/http"
	"time"

	"github.com/SIRIUS-webkit/crud-app/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)	

type RegisterUserRequest struct {
	Email string `json:"email" binding:"required,email"`
	Name string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}

// RegisterUser registers a new user
// @Summary		Registers a new user
// @Description	Registers a new user
// @Tags			auth
// @Accept			json
// @Produce		json
// @Param			user	body		registerRequest	true	"User"
// @Success		201	{object}	database.User
// @Router			/api/v1/auth/register [post]

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


type LoginUserRequest struct {
	Email string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginUserResponse struct {
	Token string `json:"token"`
}

// Login logs in a user
//
//	@Summary		Logs in a user
//	@Description	Logs in a user
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			user	body	loginRequest	true	"User"
//	@Success		200	{object}	loginResponse
//	@Router			/api/v1/auth/login [post]

func (app *application) loginUser(c *gin.Context){
	var request LoginUserRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := app.models.Users.GetByEmail(request.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "userId": user.Id,
        "exp":    time.Now().Add(time.Hour * 72).Unix(), // Token expires in 72 hours
    })

    tokenString, err := token.SignedString([]byte(app.jwtSecret))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
        return
    }

    c.JSON(http.StatusOK, LoginUserResponse{Token: tokenString})

}
	
	
	
	

