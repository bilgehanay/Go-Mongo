package main

import (
	"github.com/bilgehanay/ResponseHandler"
	"github.com/gin-gonic/gin"
)

/*#################### USER CONTROLLER ###############################*/
func createUser(c *gin.Context) {
	var user User
	r := ResponseHandler.New()
	if err := c.ShouldBindJSON(&user); err != nil {
		r.Errors = gin.H{"user": user}
		r.SendResponse(c, 3000)
		return
	}

	err := CreateUpdateUser(&user)
	if err != nil {
		r.Errors = gin.H{"user": user}
		r.SendResponse(c, 5007)
		return
	}
	r.SendResponse(c, 10000)
	return
}

func getUsers(c *gin.Context) {
	r := ResponseHandler.New()
	query := c.Request.URL.RawQuery
	users, err := GetUsers(query)
	if err != nil {
		r.Errors = gin.H{"query": query}
		r.SendResponse(c, 5007)
		return
	}
	r.Count = len(users)
	r.Data = users
	r.SendResponse(c, 10000)
	return

}

func getUser(c *gin.Context) {
	id := c.Param("id")
	r := ResponseHandler.New()
	user, err := GetUserByID(id)
	if err != nil {
		r.Errors = gin.H{"id": id}
		r.SendResponse(c, 5007)
		return
	}
	r.Data = user
	r.Count = 1
	r.SendResponse(c, 10000)
	return
}

func updateUser(c *gin.Context) {
	var user User
	r := ResponseHandler.New()
	if err := c.ShouldBindJSON(&user); err != nil {
		r.Errors = gin.H{"user": user}
		r.SendResponse(c, 3000)
		return
	}

	err := CreateUpdateUser(&user)
	if err != nil {
		r.Errors = gin.H{"user": user}
		r.SendResponse(c, 5007)
		return
	}
	r.SendResponse(c, 10000)
	return
}

func deleteUser(c *gin.Context) {
	id := c.Param("id")
	r := ResponseHandler.New()
	if err := DeleteUserByID(id); err != nil {
		r.Errors = gin.H{"id": id}
		r.SendResponse(c, 5007)
		return
	}
	r.SendResponse(c, 10000)
	return
}

/*#################### FAVORITE CONTROLLER ###########################*/
func getUserFavorites(c *gin.Context) {
	id := c.Param("uid")
	r := ResponseHandler.New()
	result, err := GetUserFavorites(id)
	if err != nil {
		r.Errors = gin.H{"id": id}
		r.SendResponse(c, 5007)
		return
	}
	r.Data = result
	r.Count = len(result)
	r.SendResponse(c, 10000)
	return
}

func putUserFavorites(c *gin.Context) {
	var fav Favorite
	r := ResponseHandler.New()
	id := c.Param("uid")
	if err := c.ShouldBindJSON(&fav); err != nil {
		r.Errors = gin.H{"id": id}
		r.SendResponse(c, 5007)
		return
	}
	if err := PutUserFavorites(id, fav); err != nil {
		r.Errors = gin.H{"id": id, "req": fav}
		r.SendResponse(c, 3000)
		return
	}
	r.SendResponse(c, 10000)
	return
}

func deleteUserFavorites(c *gin.Context) {
	uid, fid := c.Param("uid"), c.Param("fid")
	r := ResponseHandler.New()
	if err := DeleteUserFavorites(uid, fid); err != nil {
		r.Errors = gin.H{"uid": uid, "fid": fid}
		r.SendResponse(c, 3000)
		return
	}
	r.SendResponse(c, 10000)
	return
}

func updateFavorite(c *gin.Context) {
	var fav Favorite
	r := ResponseHandler.New()
	if err := c.ShouldBindJSON(&fav); err != nil {
		r.Errors = gin.H{"fav": fav}
		r.SendResponse(c, 5007)
		return
	}
	err := UpdateFavorite(fav)
	if err != nil {
		r.Errors = gin.H{"fav": fav}
		r.SendResponse(c, 3000)
		return
	}
	r.SendResponse(c, 10000)
	return
}

/*#################### COMMENT CONTROLLER ############################*/
func postComment(c *gin.Context) {
	uid := c.Param("uid")
	r := ResponseHandler.New()
	var comment map[string]interface{}

	if err := c.ShouldBindJSON(&comment); err != nil {
		r.Errors = gin.H{"id": uid}
		r.SendResponse(c, 5007)
		return
	}

	err := PostPutComment(uid, comment)
	if err != nil {
		r.Errors = gin.H{"id": uid, "req": comment}
		r.SendResponse(c, 3000)
		return
	}
	r.SendResponse(c, 10000)
	return
}

func updateComment(c *gin.Context) {
	uid := c.Param("id")
	r := ResponseHandler.New()
	var comment map[string]interface{}

	if err := c.ShouldBindJSON(&comment); err != nil {
		r.Errors = gin.H{"id": uid}
		r.SendResponse(c, 5007)
		return
	}

	err := PostPutComment(uid, comment)
	if err != nil {
		r.Errors = gin.H{"id": uid, "req": comment}
		r.SendResponse(c, 3000)
		return
	}
	r.SendResponse(c, 10000)
	return
}

func deleteComment(c *gin.Context) {
	cid := c.Param("cid")
	r := ResponseHandler.New()
	if err := DeleteComment(cid); err != nil {
		r.Errors = gin.H{"cid": cid}
		r.SendResponse(c, 3000)
		return
	}
	r.SendResponse(c, 10000)
	return
}

/*################### ORDER CONTROLLER ###############################*/
func postOrder(c *gin.Context) {
	var order Order
	r := ResponseHandler.New()
	if err := c.ShouldBindJSON(&order); err != nil {
		r.SendResponse(c, 5007)
		return
	}

	if err := CreateUpdateOrder(order); err != nil {
		r.Errors = gin.H{"req": order}
		r.SendResponse(c, 3000)
		return
	}
	r.SendResponse(c, 10000)
	return
}

func updateOrder(c *gin.Context) {
	var order Order
	r := ResponseHandler.New()
	if err := c.ShouldBindJSON(&order); err != nil {
		r.SendResponse(c, 5007)
		return
	}
	if err := CreateUpdateOrder(order); err != nil {
		r.Errors = gin.H{"req": order}
		r.SendResponse(c, 3000)
		return
	}
	r.SendResponse(c, 10000)
	return
}

func getOrders(c *gin.Context) {
	r := ResponseHandler.New()
	result, err := GetOrders()
	if err != nil {
		r.Errors = gin.H{"req": result}
		r.SendResponse(c, 5007)
		return
	}
	r.Data = result
	r.Count = len(result)
	r.SendResponse(c, 10000)
	return
}

func getUserOrders(c *gin.Context) {
	r := ResponseHandler.New()
	uid := c.Param("id")
	result, err := GetUserOrders(uid)
	if err != nil {
		r.Errors = gin.H{"req": result}
		r.SendResponse(c, 5007)
		return
	}
	r.Data = result
	r.Count = len(result)
	r.SendResponse(c, 10000)
	return
}

func deleteOrder(c *gin.Context) {
	r := ResponseHandler.New()
	oid := c.Param("id")
	if err := DeleteOrder(oid); err != nil {
		r.Errors = gin.H{"req": oid}
		r.SendResponse(c, 3000)
		return
	}
	r.SendResponse(c, 10000)
	return
}
