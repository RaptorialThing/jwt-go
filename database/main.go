package database

import (
	"context"
	"log"

	// "errors"
	"time"
	//"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/bson"
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
	Guid [36]byte `bson:"guid"`
	Refresh_token string `bson:"refresh_token"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`

}

type User struct {
	Guid [36]byte `bson: "Guid"`
	Refresh_token string ` bson: "Refresh_token"`
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

func ConvertGuid(guid string) ([36]byte, bool) {
	guidBytes := []byte(guid)
	var err bool
	var guidBytesFixed [36]byte
	err = false

	if len([36]byte{}) != len(guid) {
		err = true
		return [36]byte{}, err
	}

	for i := range [36]byte{} {
		guidBytesFixed[i] = guidBytes[i]
	}

	return guidBytesFixed, err
}

func GetUsers() ([]*UserMongo, error) {
	// filter := bson.D{
	// 	primitive.E{Key: "refresh_token", Value: "zQMDQiYCOhgHOvgSeycJPJHYNufNjJhhjUVRuSqfgqVMkPYVkURUpiFvIZRgBmyArKCtzkjkZIvaBjMkXVbWGvbq"},
	// }

	filter := bson.D{{}}

	return FilterUsers(filter)
}

func SaveUser(user User) (bool) {

	userMongo := &UserMongo{
		ID: primitive.NewObjectID(),
		Guid: user.Guid,
		Refresh_token: user.Refresh_token,
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

func UpdateUser(user User) bool {

	usersSaved,err := GetUser(user.Guid)
	userSaved := usersSaved[len(usersSaved)-1]
	id := userSaved.ID 

	if err != nil {
		return false
	}

	userSaved.Refresh_token = user.Refresh_token
	userSaved.UpdatedAt = time.Now()

	filter := bson.D{{"_id", id}}
	update := bson.D{{"$set",
	bson.D{
		{"refresh_token", userSaved.Refresh_token},
		{"updated_at", time.Now()},
	},
	}}

	_, errUpd := usersCollection.UpdateOne(
		ctx,
		filter,
        update,
	)

	usrYX,_ := GetUser(user.Guid)

	if usrYX != nil {

	}
	// fmt.Println(usrYX[0].Refresh_token)

	if errUpd != nil {
		return false
	}

	return true	
}

func GetUser(guid [36]byte) ([]*UserMongo, error) {
	
	filter := bson.D{
		primitive.E{Key:"guid", Value: guid},
	}
	result, err := FilterUsers(filter)

	return result,err
}

func main() {

}
