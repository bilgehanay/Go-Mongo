package main

import (
	"fmt"
	"github.com/ajclopez/mgs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Primitives struct{}

func (primitives *Primitives) ObjectID(oidStr string) (interface{}, error) {
	return primitive.ObjectIDFromHex(oidStr)
}

func CreateUpdateUser(user *User) error {
	opt := options.Update().SetUpsert(true)
	if user.ID.IsZero() {

		user.ID = primitive.NewObjectID()
	}
	filter := bson.M{"_id": user.ID}
	_, err := db.UpdateOne(ctx, filter, bson.M{"$set": user}, opt)
	return err
}

// GetUsers Query filter added, if the query is empty it returns all
/*
$eq			key=val				type=public
$ne			key!=val			status!=SENT
$gt			key>val				price>5
$gte		key>=val			price>=9
$lt			key<val				date<2020-01-01T14:00:00.000Z
$lte		key<=val			priority<=-5
$in			key=val1,val2		status=QUEUED,DEQUEUED
$nin		key!=val1,val2		status!=QUEUED,DEQUEUED
$exists		key					email
$exists		!key				!email
$regex		key=/value/<opts>	email=/@gmail\.com$/
$regex		key!=/value/<opts>	phone!=/^58/
*/
func GetUsers(query string) ([]*User, error) {
	var users []*User
	opt := mgs.FindOption()
	opt.SetMaxLimit(100)
	queryHandling := mgs.NewQueryHandler(&Primitives{})
	result, err := queryHandling.MongoGoSearch(query, opt)

	if err != nil {
		return nil, err
	}

	findOpts := options.Find()
	findOpts.SetLimit(result.Limit)
	findOpts.SetSkip(result.Skip)
	findOpts.SetSort(result.Sort)
	findOpts.SetProjection(result.Projection)

	cur, err := db.Find(ctx, result.Filter, findOpts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var user User
		err := cur.Decode(&user)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

func GetUserByID(id string) (*User, error) {
	user_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var user User
	filter := bson.M{"_id": user_id}
	err = db.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func DeleteUserByID(id string) error {
	user_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": user_id}
	_, err = db.DeleteOne(ctx, filter)
	return err
}

func GetUserFavorites(id string) ([]Favorite, error) {
	var user User

	user_id, err := primitive.ObjectIDFromHex(id)

	filter := bson.M{"_id": user_id}

	err = db.FindOne(ctx, filter).Decode(&user)
	fmt.Println(user.Favorites)
	if err != nil {
		return nil, err
	}
	return user.Favorites, nil
}

func PutUserFavorites(id string, fav Favorite) error {
	user_id, err := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": user_id}
	update := bson.M{
		"$push": bson.M{
			"favorites": fav,
		},
	}
	_, err = db.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func DeleteUserFavorites(uid, pid string) error {
	user_id, err := primitive.ObjectIDFromHex(uid)
	product_id, err := primitive.ObjectIDFromHex(pid)
	filter := bson.M{"_id": user_id}
	update := bson.M{
		"$pull": bson.M{
			"favorites": bson.M{
				"product_id": product_id,
			},
		},
	}
	_, err = db.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}
