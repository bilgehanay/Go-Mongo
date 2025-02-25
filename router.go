package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

var (
	HttpServer *http.Server
	Corsconfig cors.Config
)

func init() {
	Corsconfig = cors.DefaultConfig()
	Corsconfig.AllowAllOrigins = true
	Corsconfig.AllowHeaders = []string{"*"}
	Corsconfig.AllowMethods = []string{"GET", "POST", "PUT", "OPTIONS", "DELETE"}

	HttpServer = &http.Server{
		Addr:         config.Port,
		Handler:      router(),
		ReadTimeout:  50 * time.Second,
		WriteTimeout: 100 * time.Second,
	}
}

func router() http.Handler {
	fmt.Println("Api starting on " + config.Port)
	r := gin.New()
	r.Use(gin.Recovery())
	gin.SetMode(gin.DebugMode)
	r.Use(cors.New(Corsconfig))
	public := r.Group("/api")
	{
		public.POST("", createUser)
		public.GET("/", getUsers)
		public.GET("/:id", getUser)
		public.PUT("", updateUser)
		public.DELETE("/:id", deleteUser)

		public.GET("/favorite/:uid", getUserFavorites)
		public.PUT("/favorite/:uid", putUserFavorites)
		public.DELETE("/favorite/:uid/:fid", deleteUserFavorites)

		public.PUT("/favorite/update", updateFavorite)

		public.POST("/comment/:uid", postComment)
		public.PUT("/comment/:uid", updateComment)
		public.DELETE("/comment/:cid", deleteComment)

		public.POST("/order", postOrder)
		public.PUT("/order", updateOrder)
		public.GET("/orders", getOrders)
		public.GET("/order/:id", getUserOrders)
		public.DELETE("/order/:id", deleteOrder)
	}
	return r
}
