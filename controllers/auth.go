package controllers

import (
	"jwt-authorizer/models"
	"jwt-authorizer/utils/token"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RegisterInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u := models.User{}

	u.Username = input.Username
	u.Password = input.Password

	_, err := u.SaveUser()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "registration success!"})
}

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	user := models.User{}
	
	user.Username = input.Username
	user.Password = input.Password

	
	token, err := models.LoginCheck(user.Username, user.Password)


	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username or password incorrect"})
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func CurrentUser(context *gin.Context){
	user_id, err := token.ExtractTokenID(context)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err})
	}

	user, err := models.GetUserById(user_id)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message":"success", "data": user})
}
