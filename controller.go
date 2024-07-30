package main

import (
	"GO-MONGO/ResponseHandler"
	"github.com/gin-gonic/gin"
)

/*#################### USER CONTROLLER ###############################*/
func createUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		ResponseHandler.HandleError(c, 5007, "", nil, gin.H{})
		return
	}

	err := CreateUpdateUser(&user)
	if err != nil {
		ResponseHandler.HandleError(c, 5007, "", nil, gin.H{"req": user})
		return
	}
	ResponseHandler.HandleSuccess(c, nil, 0)
	return
}

func getUsers(c *gin.Context) {
	query := c.Request.URL.RawQuery
	users, err := GetUsers(query)
	if err != nil {
		ResponseHandler.HandleError(c, 5007, "", nil, gin.H{"query": query})
		return
	}
	ResponseHandler.HandleSuccess(c, users, len(users))
	return

}

func getUser(c *gin.Context) {
	id := c.Param("id")
	user, err := GetUserByID(id)
	if err != nil {
		ResponseHandler.HandleError(c, 5007, "", nil, gin.H{"id": id})
		return
	}
	ResponseHandler.HandleSuccess(c, user, 1)
	return
}

func updateUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		ResponseHandler.HandleError(c, 5007, "", nil, gin.H{})
		return
	}

	err := CreateUpdateUser(&user)
	if err != nil {
		ResponseHandler.HandleError(c, 3000, "", nil, gin.H{"req": user})
		return
	}
	ResponseHandler.HandleSuccess(c, nil, 0)
	return
}

func deleteUser(c *gin.Context) {
	id := c.Param("id")
	if err := DeleteUserByID(id); err != nil {
		ResponseHandler.HandleError(c, 3000, "", nil, gin.H{"id": id})
		return
	}
	ResponseHandler.HandleSuccess(c, nil, 0)
	return
}

/*#################### FAVORITE CONTROLLER ###########################*/
func getUserFavorites(c *gin.Context) {
	id := c.Param("uid")
	result, err := GetUserFavorites(id)
	if err != nil {
		ResponseHandler.HandleError(c, 3000, "", nil, gin.H{"id": id})
		return
	}
	ResponseHandler.HandleSuccess(c, result, len(result))
	return
}

func putUserFavorites(c *gin.Context) {
	var fav Favorite
	id := c.Param("uid")
	if err := c.ShouldBindJSON(&fav); err != nil {
		ResponseHandler.HandleError(c, 5007, "", nil, gin.H{"id": id})
		return
	}
	if err := PutUserFavorites(id, fav); err != nil {
		ResponseHandler.HandleError(c, 3000, "", nil, gin.H{"id": id, "req": fav})
		return
	}
	ResponseHandler.HandleSuccess(c, nil, 0)
	return
}

func deleteUserFavorites(c *gin.Context) {
	uid, fid := c.Param("uid"), c.Param("fid")
	if err := DeleteUserFavorites(uid, fid); err != nil {
		ResponseHandler.HandleError(c, 3000, "", nil, gin.H{"fid": fid, "uid": uid})
		return
	}
	ResponseHandler.HandleSuccess(c, nil, 0)
	return
}

func updateFavorite(c *gin.Context) {
	var fav Favorite
	if err := c.ShouldBindJSON(&fav); err != nil {
		ResponseHandler.HandleError(c, 5007, "", nil, gin.H{})
		return
	}
	err := UpdateFavorite(fav)
	if err != nil {
		ResponseHandler.HandleError(c, 3000, "", nil, gin.H{"req": fav})
		return
	}
	ResponseHandler.HandleSuccess(c, nil, 0)
	return
}

/*#################### COMMENT CONTROLLER ############################*/
func postComment(c *gin.Context) {
	uid := c.Param("uid")
	var comment map[string]interface{}

	if err := c.ShouldBindJSON(&comment); err != nil {
		ResponseHandler.HandleError(c, 5007, "", nil, gin.H{"id": uid})
		return
	}

	err := PostPutComment(uid, comment)
	if err != nil {
		ResponseHandler.HandleError(c, 3000, "", nil, gin.H{"id": uid, "req": comment})
		return
	}
	ResponseHandler.HandleSuccess(c, nil, 0)
	return
}

func updateComment(c *gin.Context) {
	uid := c.Param("id")
	var comment map[string]interface{}

	if err := c.ShouldBindJSON(&comment); err != nil {
		ResponseHandler.HandleError(c, 5007, "", nil, gin.H{"id": uid})
		return
	}

	err := PostPutComment(uid, comment)
	if err != nil {
		ResponseHandler.HandleError(c, 3000, "", nil, gin.H{"id": uid, "req": comment})
		return
	}
	ResponseHandler.HandleSuccess(c, nil, 0)
	return
}

func deleteComment(c *gin.Context) {
	cid := c.Param("cid")
	if err := DeleteComment(cid); err != nil {
		ResponseHandler.HandleError(c, 3000, "", nil, gin.H{"id": cid})
		return
	}
	ResponseHandler.HandleSuccess(c, nil, 0)
	return
}

/*################### ORDER CONTROLLER ###############################*/
func postOrder(c *gin.Context) {
	var order Order

	if err := c.ShouldBindJSON(&order); err != nil {
		ResponseHandler.HandleError(c, 5007, "", nil, gin.H{})
		return
	}

	if err := CreateUpdateOrder(order); err != nil {
		ResponseHandler.HandleError(c, 3000, "", nil, gin.H{"req": order})
		return
	}
	ResponseHandler.HandleSuccess(c, nil, 0)
	return
}

func updateOrder(c *gin.Context) {
	var order Order
	if err := c.ShouldBindJSON(&order); err != nil {
		ResponseHandler.HandleError(c, 5007, "", nil, gin.H{})
		return
	}
	if err := CreateUpdateOrder(order); err != nil {
		ResponseHandler.HandleError(c, 3000, "", nil, gin.H{"req": order})
		return
	}
	ResponseHandler.HandleSuccess(c, nil, 0)
	return
}

func getOrders(c *gin.Context) {
	result, err := GetOrders()
	if err != nil {
		ResponseHandler.HandleError(c, 3000, "", nil, gin.H{"req": result})
		return
	}
	ResponseHandler.HandleSuccess(c, result, len(result))
	return
}

func getUserOrders(c *gin.Context) {
	uid := c.Param("id")
	result, err := GetUserOrders(uid)
	if err != nil {
		ResponseHandler.HandleError(c, 3000, "", nil, gin.H{"req": result})
		return
	}
	ResponseHandler.HandleSuccess(c, result, len(result))
	return
}

func deleteOrder(c *gin.Context) {
	oid := c.Param("id")
	if err := DeleteOrder(oid); err != nil {
		ResponseHandler.HandleError(c, 3000, "", nil, gin.H{"req": oid})
		return
	}
	ResponseHandler.HandleSuccess(c, nil, 0)
	return
}
