package controllers

import (
	"api/api/auth"
	"api/api/security"
	"api/api/utils/formaterror"
	"api/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func (server *Server) Login(c *gin.Context) {
	errList = map[string]string{}
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  "Unable to get request",
		})
		return
	}
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  "Can not unmarshal body",
		})
		return
	}
	user.Prepare()
	errorMessages := user.Validate("login")
	if len(errorMessages) > 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errorMessages,
		})
		return
	}
	userData, err := server.SignIn(user.Email, user.Password)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  formattedError,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": userData,
	})
}

func (server *Server) SignIn(email, password string) (map[string]interface{}, error) {
	var err error
	userData := make(map[string]interface{})
	user := models.User{}
	err = server.DB.Debug().Model(models.User{}).Where("email = ?", email).Take(&user).Error
	if err != nil {
		fmt.Printf("there is error while getting user %v", err)
		return nil, err
	}
	err = security.VerifyPassword(user.Password, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		fmt.Printf("Error while comparing the password %v", err)
		return nil, err
	}
	token, err := auth.CreateToken(user.ID)
	if err != nil {
		fmt.Printf("error while creating the token %v", err)
		return nil, err
	}
	userData["token"] = token
	userData["id"] = user.ID
	userData["firstname"] = user.FirstName
	userData["lastname"] = user.LastName
	userData["email"] = user.Email
	userData["avatar"] = user.Avatar
	userData["username"] = user.Username

	return userData, nil
}
