package main

import (
	"errors"
	"github.com/ajclopez/mgs"
	"github.com/goccy/go-json"
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
	result, err := db.UpdateOne(ctx, filter, bson.M{"$set": user}, opt)
	if result != nil && result.MatchedCount == 0 {
		return errors.New("user does not exist")
	}
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
	if len(users) == 0 {
		return nil, errors.New("no users found")
	}
	return users, nil
}

func GetUserByID(id string) (*User, error) {
	userId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var user User
	filter := bson.M{"_id": userId}
	err = db.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func DeleteUserByID(id string) error {
	userId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": userId}
	result, err := db.DeleteOne(ctx, filter)
	if result != nil && result.DeletedCount == 0 {
		return errors.New("user does not exist")
	}
	return err
}

func GetUserFavorites(id string) ([]Favorite, error) {
	var user User

	userId, err := primitive.ObjectIDFromHex(id)

	filter := bson.M{"_id": userId}

	err = db.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return user.Favorites, nil
}

func PutUserFavorites(id string, fav Favorite) error {
	userId, err := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": userId}
	update := bson.M{
		"$push": bson.M{
			"favorites": fav,
		},
	}
	result, err := db.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.ModifiedCount == 0 {
		return errors.New("favorite not found")
	}

	return nil
}

// Arrayden çıkarılması gereken objeyi sahip user id si ve obje id si ile buluyor
func DeleteUserFavorites(uid, pid string) error {
	userId, err := primitive.ObjectIDFromHex(uid)
	productId, err := primitive.ObjectIDFromHex(pid)
	filter := bson.M{"_id": userId}
	update := bson.M{
		"$pull": bson.M{
			"favorites": bson.M{
				"product_id": productId,
			},
		},
	}
	result, err := db.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("favorite not found")
	}
	return nil
}

// Favorite field changes applying for every user favorites array
func UpdateFavorite(fav Favorite) error {
	filter := bson.M{"favorites.product_id": fav.ProductID}

	var fields map[string]interface{}
	data, _ := json.Marshal(fav)
	json.Unmarshal(data, &fields)

	updateData := bson.M{}
	for k, v := range fields {
		if k != "product_id" {
			updateData["favorites.$[p]."+k] = v
		}
	}

	update := bson.M{"$set": updateData}
	arrayFilters := options.Update().SetArrayFilters(options.ArrayFilters{
		Filters: []interface{}{bson.M{"p.product_id": fav.ProductID}},
	})
	result, err := db.UpdateMany(ctx, filter, update, arrayFilters)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("favorite not found")
	}
	return nil
}

func PostPutComment(uid string, comment map[string]interface{}) error {
	userId, err := primitive.ObjectIDFromHex(uid)
	if err != nil {
		return err
	}

	if _, ok := comment["cid"]; !ok {
		comment["cid"] = primitive.NewObjectID()
	} else {
		comment["cid"], _ = primitive.ObjectIDFromHex(comment["cid"].(string))
	}

	filter := bson.M{"_id": userId, "comments.cid": comment["cid"]}
	update := bson.M{"$set": bson.M{
		"comments.$[c]": comment,
	}}
	arrayFilters := options.Update().SetArrayFilters(options.ArrayFilters{
		Filters: []interface{}{bson.M{"c.cid": comment["cid"]}},
	})
	result, err := db.UpdateOne(ctx, filter, update, arrayFilters)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		insert := bson.M{
			"$push": bson.M{
				"comments": comment,
			},
		}
		_, err = db.UpdateOne(ctx, bson.M{"_id": userId}, insert)
	}
	return err
}

// Arrayden çıkarılması gereken objeyi direkt olarak obje id si ile buluyor
func DeleteComment(cid string) error {
	commentId, err := primitive.ObjectIDFromHex(cid)
	filter := bson.M{"comments.cid": commentId}
	update := bson.M{
		"$pull": bson.M{
			"comments": bson.M{
				"cid": commentId,
			},
		},
	}
	result, err := db.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("comment not found")
	}
	return err
}
