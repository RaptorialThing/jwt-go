package main

import (
	"context"
	"log"

	// "errors"
	"time"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/bson"
	"example.com/jwtgo"
	// "github.com/urfave/cli/v2"
)

var usersCollection *mongo.Collection
var ctx = context.TODO()

func init() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	usersCollection = client.Database("test").Collection("users")
} 

type UserMongo struct {
	ID        primitive.ObjectID `bson:"_id"`
	guid [36]byte `bson: "guid"`
	refresh_token string `bson: "refresh_token"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`

}

type User struct {
	guid [36]byte `bson: "GUID"`
	refresh_token string ` bson: "refresh_token"`
}

func createUser(user *UserMongo) error {
	_, err := usersCollection.InsertOne(ctx, user)
  return err
}

func getAll() ([]*UserMongo, error) {
	  filter := bson.D{{}}
	  return FilterUsers(filter)
  }

func FilterUsers(filter interface{}) ([]*UserMongo, error) {

	var users []*UserMongo

	cur, err := usersCollection.Find(ctx, filter)
	if err != nil {
		return users, err
	}

	for cur.Next(ctx) {
		var t UserMongo
		err := cur.Decode(&t)
		if err != nil {
			return users, err
		}

		users = append(users, &t)
	}

	if err := cur.Err(); err != nil {
		return users, err
	}

	cur.Close(ctx)

	if len(users) == 0 {
		return users, mongo.ErrNoDocuments
	}

	return users, nil
}

func GetUsers() ([]*UserMongo, error) {
	// filter := bson.D{
	// 	primitive.E{Key: "refresh_token", Value: "zQMDQiYCOhgHOvgSeycJPJHYNufNjJhhjUVRuSqfgqVMkPYVkURUpiFvIZRgBmyArKCtzkjkZIvaBjMkXVbWGvbq"},
	// }

	filter := bson.D{{}}

	return FilterUsers(filter)
}

func SaveUser(user User) (bool) {

	userMongo := UserMongo{
		ID: primitive.NewObjectID(),
		guid: user.guid,
		refresh_token: user.refresh_token,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	userResult, err := usersCollection.InsertOne(ctx, userMongo)

	if err != nil {
		return false
	}
	if userResult != nil {

	}
	return true

}

func GetUser(guid [36]byte) ([]*UserMongo, error) {
	// client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017/"))
    // if err != nil {
    //     log.Fatal(err)
    // }
	// ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
    // err = client.Connect(ctx)
    // if err != nil {
    //     log.Fatal(err)
    // }
    // defer client.Disconnect(ctx)

	// db := client.Database("test")
	// usersCollection := db.Collection("users")

	// var result UserMongo
	// err = usersCollection.FindOne(context.TODO(),bson.D{primitive.E{Key:"guid",Value:guid}}).Decode(&result)
	// if err != nil {
	// 	if err == mongo.ErrNoDocuments {
	// 		return UserMongo{},err
	// 	}
	// } 
	filter := bson.D{
		primitive.E{Key:"GUID", Value: guid},
	}
	result, err := FilterUsers(filter)

	return result,err
}

func main() {
	// byteGuid,errConv := jwtgo.ConvertGuid("e8f39331-bc2e-4392-97b1-2328b3c63ab6")
	// if errConv {
	// 	fmt.Println("Error conv")
	// }


	// userToSave := User{
	// 	guid: byteGuid,
	// 	refresh_token: "$2a$08$h5ZQnwD2GI5C9eDG3ECMcOaMylqkk12Oem.WgloitHcV0nmqDZvGm",
	// }

	// resSaveUsr := SaveUser(userToSave)
	// if resSaveUsr {
	// 	fmt.Println("Save success")
	// } else {
	// 	fmt.Println("Save error")
	// }



	byteGuidTwo,errConv := jwtgo.ConvertGuid("e8f39331-bc2e-4392-97b1-2328b3c63ab6")

	if errConv {
		fmt.Println("Error convert")
	}

	usr,err := GetUser(byteGuidTwo)
	if err != nil {
		fmt.Println("Error get User "+err.Error())
	}

	fmt.Println(usr)

}
