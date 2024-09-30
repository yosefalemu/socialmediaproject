package controllers

import (
	"api/api/auth"
	"api/api/utils/formaterror"
	"api/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (server *Server) CreateUser(c *gin.Context) {
	fmt.Println("USER CREATION STARTED")
	errList = map[string]string{}
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		errList["Invalid_body"] = "Unable to get request"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		errList["Unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}
	err = user.BeforeSave()
	if err != nil {
		fmt.Println("error while hashing password", err)
	}
	user.Prepare()
	errorMessages := user.Validate("create")
	if len(errorMessages) > 0 {
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}
	userCreated, err := user.SaveUser(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		errList = formattedError
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  errList,
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"status":   http.StatusCreated,
		"response": userCreated,
	})
}

// GET ALL THE USERS IN THE SYSTEM
func (server *Server) GetUsers(c *gin.Context) {
	errList = map[string]string{}
	user := models.User{}
	users, err := user.FindAllUsers(server.DB)
	if err != nil {
		errList["No_user"] = "No user found"
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  errList,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": users,
	})
}

// GET SINGLE USER FROM THE SYSTEM
func (server *Server) GetUser(c *gin.Context) {
	errList = map[string]string{}
	userId := c.Param("id")
	uid, err := strconv.ParseInt(userId, 10, 32)
	if err != nil {
		errList["invalid_request"] = "Invalid Request"
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  errList,
		})
		return
	}
	user := models.User{}
	userGotten, err := user.FindUserById(server.DB, uint32(uid))
	if err != nil {
		errList["No_User"] = "No User Found"
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  errList,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": userGotten,
	})
}

// UPDATE SINGLE USER
func (server *Server) UpdateUser(c *gin.Context) {
	errList = map[string]string{}
	userId := c.Param("id")
	uid, err := strconv.ParseUint(userId, 10, 32)
	if err != nil {
		errList["invalid_request"] = "Invalid Request"
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  errList,
		})
		return
	}
	//GET USER ID FROM THE TOKEN
	tokenId, err := auth.ExtractTokenID(c.Request)
	fmt.Println("TOKEN ID FOUND", tokenId)
	if err != nil {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}
	if tokenId != 0 && tokenId != uint32(uid) {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}
	//AFTER THIS CHECK START PROCESSING THE REQUEST
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		errList["Unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}
	requestBody := map[string]string{}
	err = json.Unmarshal(body, &requestBody)
	if err != nil {
		errList["Unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}
	//CHECK FOR THE PREVIOUS DETAILS
	formerUser := models.User{}
	err = server.DB.Debug().Model(models.User{}).Where("id = ?", uid).Take(&formerUser).Error
	if err != nil {
		errList["User_invalid"] = "The user does not exist"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}
	fmt.Println("Requested body", requestBody)
	if newFirstName, exists := requestBody["first_name"]; exists {
		formerUser.FirstName = newFirstName
	}
	if newLastName, exists := requestBody["last_name"]; exists {
		formerUser.LastName = newLastName
	}
	if newEmail, exists := requestBody["email"]; exists {
		formerUser.Email = newEmail
	}
	if newAvatar, exists := requestBody["avatar"]; exists {
		formerUser.Avatar = newAvatar
	}

	formerUser.Prepare()
	errorMessages := formerUser.Validate("update")
	if len(errorMessages) > 0 {
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
	}
	updatedUser, err := formerUser.UpdateUser(server.DB, uint32(uid))
	if err != nil {
		errList := formaterror.FormatError(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  errList,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": updatedUser,
	})
}
