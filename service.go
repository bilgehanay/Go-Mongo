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

/* ############# USER SERVICE ############# */
func CreateUpdateUser(user *User) error {
	opt := options.Update().SetUpsert(true)
	if user.ID.IsZero() {
		user.ID = primitive.NewObjectID()
	}
	filter := bson.M{"_id": user.ID}
	_, err := userdb.UpdateOne(ctx, filter, bson.M{"$set": user}, opt)
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

	cur, err := userdb.Find(ctx, result.Filter, findOpts)
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
	err = userdb.FindOne(ctx, filter).Decode(&user)
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
	result, err := userdb.DeleteOne(ctx, filter)
	if result != nil && result.DeletedCount == 0 {
		return errors.New("user does not exist")
	}
	return err
}

/* ############# FAVORITE SERVICE ############# */
func GetUserFavorites(id string) ([]Favorite, error) {
	var user User

	userId, err := primitive.ObjectIDFromHex(id)

	filter := bson.M{"_id": userId}

	err = userdb.FindOne(ctx, filter).Decode(&user)
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
	result, err := userdb.UpdateOne(ctx, filter, update)
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
	result, err := userdb.UpdateOne(ctx, filter, update)
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
	err := json.Unmarshal(data, &fields)
	if err != nil {
		return err
	}

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
	result, err := userdb.UpdateMany(ctx, filter, update, arrayFilters)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("favorite not found")
	}
	return nil
}

/* ############# COMMENT SERVICE ############# */
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
	result, err := userdb.UpdateOne(ctx, filter, update, arrayFilters)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		insert := bson.M{
			"$push": bson.M{
				"comments": comment,
			},
		}
		_, err = userdb.UpdateOne(ctx, bson.M{"_id": userId}, insert)
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
	result, err := userdb.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("comment not found")
	}
	return err
}

/* ############# ORDER SERVICE ############# */
func CreateUpdateOrder(o Order) error {
	opt := options.Update().SetUpsert(true)
	if o.ID.IsZero() {
		o = NewOrder(primitive.NewObjectID(), o.UserID, o.Name, o.Quantity, o.Price)
	}
	filter := bson.M{"_id": o.ID}
	_, err := orderdb.UpdateOne(ctx, filter, bson.M{"$set": o}, opt)
	return err
}

type Result struct {
	User   User    `json:"user"`
	Orders []Order `json:"orders"`
}

func GetOrders() ([]Result, error) {
	var result []Result
	pipeline := bson.A{
		bson.D{{"$lookup", bson.D{
			{"from", "order"},
			{"localField", "_id"},
			{"foreignField", "user_id"},
			{"as", "orders"},
		}}},
	}
	cur, err := userdb.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var rawResults []bson.M
	if err = cur.All(ctx, &rawResults); err != nil {
		return nil, err
	}

	for _, rawResult := range rawResults {
		var user User
		var orders []Order

		bsonBytes, _ := bson.Marshal(rawResult)
		bson.Unmarshal(bsonBytes, &user)

		if rawOrders, ok := rawResult["orders"]; ok {
			for _, rawOrder := range rawOrders.(primitive.A) {
				var order Order
				bsonBytes, _ := bson.Marshal(rawOrder)
				bson.Unmarshal(bsonBytes, &order)
				orders = append(orders, order)
			}
		}

		result = append(result, Result{
			User:   user,
			Orders: orders,
		})
	}

	if len(result) == 0 {
		return nil, errors.New("no orders found")
	}
	return result, nil
}

func GetUserOrders(uid string) ([]Order, error) {
	userId, err := primitive.ObjectIDFromHex(uid)
	if err != nil {
		return nil, err
	}
	var orders []Order
	filter := bson.M{"user_id": userId}
	cur, err := orderdb.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var order Order
		err := cur.Decode(&order)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func DeleteOrder(oid string) error {
	filter := bson.M{"_id": oid}
	result, err := orderdb.DeleteOne(ctx, filter)
	if result != nil && result.DeletedCount == 0 {
		return errors.New("order not found")
	}
	return err
}
