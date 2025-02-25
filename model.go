package main

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type OrderStatus string

const (
	Pending   OrderStatus = "pending"
	Completed OrderStatus = "completed"
	Shipped   OrderStatus = "shipped"
	Cancelled OrderStatus = "cancelled"
)

type MongoConfig struct {
	ConnectionName   string `json:"connectionName"`
	ConnectionString string `json:"connectionString"`
	Collection       map[string]struct {
		N string `json:"n"` // name
		D string `json:"d"` // db
		C string `json:"c"` // col
	} `json:"collection"`
}

type Address struct {
	Street  string `json:"street" bson:"street,omitempty"`
	State   string `json:"state" bson:"state,omitempty"`
	City    string `json:"city" bson:"city,omitempty"`
	ZipCode string `json:"zip_code" bson:"zip_code,omitempty"`
	Country string `json:"country" bson:"country,omitempty"`
}

type Favorite struct {
	ProductID       primitive.ObjectID `json:"product_id,omitempty" bson:"product_id,omitempty"`
	ProductName     string             `json:"product_name,omitempty" bson:"product_name,omitempty"`
	ProductCategory int                `json:"product_category,omitempty" bson:"product_category,omitempty"`
}

type User struct {
	ID        primitive.ObjectID       `json:"id,omitempty," bson:"_id,omitempty"`
	Name      string                   `json:"name" bson:"name" validate:"required,min=2,max=32"`
	Surname   string                   `json:"surname" bson:"surname" validate:"required,min=2,max=32"`
	Email     string                   `json:"email" bson:"email" validate:"required,email"`
	Password  string                   `json:"password" bson:"password" validate:"required,min=6,max=32"`
	Age       int                      `json:"age" bson:"age" validate:"required,min=18,max=120"`
	Address   Address                  `json:"address" bson:"address,omitempty" validate:"required"`
	Favorites []Favorite               `json:"favorites" bson:"favorites,omitempty"`
	Comments  []map[string]interface{} `json:"comments" bson:"comments,omitempty"`
}

type Order struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserID    primitive.ObjectID `json:"user_id,omitempty" bson:"user_id,omitempty"`
	Name      string             `json:"name" bson:"name" validate:"required,min=2,max=32"`
	Quantity  int                `json:"quantity" bson:"quantity" validate:"required,min=1,max=100"`
	OrderDate primitive.DateTime `json:"order_date" bson:"order_date" validate:"required"`
	Price     float64            `json:"price" bson:"price" validate:"required,min=0"`
	Status    OrderStatus        `json:"status" bson:"status" validate:"required"`
}

func NewOrder(orderID, userID primitive.ObjectID, name string, quantity int, price float64) Order {
	return Order{
		ID:        orderID,
		UserID:    userID,
		Name:      name,
		Quantity:  quantity,
		OrderDate: primitive.NewDateTimeFromTime(time.Now()),
		Price:     price,
		Status:    Pending,
	}
}
