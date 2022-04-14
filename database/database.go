package database

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
	// "github.com/urfave/cli/v2"
)

var collection *mongo.Collection
var ctx = context.TODO()

func init() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	collection = client.Database("test").Collection("users")
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
	_, err := collection.InsertOne(ctx, user)
  return err
}

func getAll() ([]*UserMongo, error) {
	  filter := bson.D{{}}
	  return filterUsers(filter)
  }

  func filterUsers(filter interface{}) ([]*UserMongo, error) {
	var users []*UserMongo

	cur, err := collection.Find(ctx, filter)
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

func getUsers(guid [36]byte) ([]*UserMongo, error) {
	filter := bson.D{
		primitive.E{Key: "refresh_token", Value: "zQMDQiYCOhgHOvgSeycJPJHYNufNjJhhjUVRuSqfgqVMkPYVkURUpiFvIZRgBmyArKCtzkjkZIvaBjMkXVbWGvbq"},
	}

	return filterUsers(filter)
}

func saveUser(user User) (bool) {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017/"))
    if err != nil {
        log.Fatal(err)
    }
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
    err = client.Connect(ctx)
    if err != nil {
        log.Fatal(err)
    }
    defer client.Disconnect(ctx)

	db := client.Database("test")
	usersCollection := db.Collection("users")

	userMongo := UserMongo{
		guid: user.guid,
		refresh_token: user.refresh_token,
	}

	userResult, err := usersCollection.InsertOne(ctx, userMongo)

	if err != nil {
		return false
	}
	if userResult != nil {

	}
	return true

}

func getUser(guid [36]byte) (UserMongo, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017/"))
    if err != nil {
        log.Fatal(err)
    }
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
    err = client.Connect(ctx)
    if err != nil {
        log.Fatal(err)
    }
    defer client.Disconnect(ctx)

	db := client.Database("test")
	usersCollection := db.Collection("users")

	var result UserMongo
	err = usersCollection.FindOne(context.TODO(),bson.D{primitive.E{Key:"guid",Value:guid}}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return UserMongo{},err
		}
	} 

	return result,err
}

func main() {

	user := &UserMongo{
		ID:        primitive.NewObjectID(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		guid: [36]byte{},
		refresh_token: "zQMDQiYCOhgHOvgSeycJPJHYNufNjJhhjUVRuSqfgqVMkPYVkURUpiFvIZRgBmyArKCtzkjkZIvaBjMkXVbWGvbq",
	}

	var gUID [36]byte
	for i:=0;i<35;i++{
		gUID[i] = byte(i)
	} 

	usrSave := User{
		guid: gUID,
		refresh_token: "123",
	}

	resSaveUsr := saveUser(usrSave)
	if resSaveUsr {
		fmt.Println("Save success")
	} else {
		fmt.Println("Save error")
	}

	if user != nil {

	}

	usr,err := getUser(gUID)
	if err != nil {
		fmt.Println("Error "+err.Error())
	}

	fmt.Println(usr)

}
