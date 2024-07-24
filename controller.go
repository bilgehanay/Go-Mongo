package main

import (
	"github.com/gin-gonic/gin"
	"go/types"
	"net/http"
)

func createRes(success bool, message string, data interface{}) gin.H {
	if success {
		return gin.H{"success": success, "message": message, "data": data}
	}
	return gin.H{"success": success, "data": message}
}

func createUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, createRes(false, err.Error(), nil))
		return
	}

	err := CreateUpdateUser(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, createRes(false, err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, createRes(true, "User Is Created", types.Struct{}))
	return
}

func getUsers(c *gin.Context) {
	query := c.Request.URL.RawQuery
	users, err := GetUsers(query)
	if err != nil {
		c.JSON(http.StatusBadRequest, createRes(false, err.Error(), types.Struct{}))
		return
	}
	c.JSON(http.StatusOK, createRes(true, "", users))
	return

}

func getUser(c *gin.Context) {
	id := c.Param("id")
	user, err := GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, createRes(false, err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, createRes(true, "", user))
	return
}

func updateUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, createRes(false, err.Error(), nil))
		return
	}

	err := CreateUpdateUser(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, createRes(false, err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, createRes(true, "User Data Is Saved", nil))
	return
}

func deleteUser(c *gin.Context) {
	id := c.Param("id")
	if err := DeleteUserByID(id); err != nil {
		c.JSON(http.StatusBadRequest, createRes(false, err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, createRes(true, "User Is Deleted", nil))
	return
}
