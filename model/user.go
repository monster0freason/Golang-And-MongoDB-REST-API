// package models

// import "gopkg.in/mgo.v2/bson"

// // User represents the user model
// type User struct {
// 	ID     bson.ObjectId `json:"id" bson:"_id"`
// 	Name   string        `json:"name" bson:"name"`
// 	Gender string        `json:"gender" bson:"gender"`
// 	Age    int           `json:"age" bson:"age"`
// }

package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// User represents the user model
type User struct {
	ID     primitive.ObjectID `json:"id" bson:"_id"`
	Name   string             `json:"name" bson:"name"`
	Gender string             `json:"gender" bson:"gender"`
	Age    int                `json:"age" bson:"age"`
}

