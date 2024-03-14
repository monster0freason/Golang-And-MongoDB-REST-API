// package controllers

// import (
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"gopkg.in/mgo.v2"
// 	"github.com/julienschmidt/httprouter"
// 	"gopkg.in/mgo.v2/bson"
// 	"github.com/monster0freason/Golang-And-MongoDB-REST-API/models"
// )

// // UserController represents the controller for handling user-related requests
// type UserController struct {
// 	session *mgo.Session
// }

// // NewUserController creates a new instance of UserController
// func NewUserController(s *mgo.Session) *UserController {
// 	return &UserController{s}
// }

// // GetUser retrieves a user by ID
// func (uc UserController) GetUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
// 	id := p.ByName("id")

// 	if !bson.IsObjectIdHex(id) {
// 		w.WriteHeader(http.StatusNotFound)
// 	}

// 	oid := bson.ObjectIdHex(id)

// 	u := models.User{}

// 	if err := uc.session.DB("mongogolang").C("users").FindId(oid).One(&u); err != nil {
// 		w.WriteHeader(http.StatusNotFound)
// 		return
// 	}

// 	uJSON, err := json.Marshal(u)
// 	if err != nil {
// 		fmt.Println(err)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	fmt.Fprintf(w, "%s\n", uJSON)
// }

// func (uc UserController) CreateUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
//     // Initialize a variable u of type models.User
//     var u models.User

//     // Decode the JSON data from the request body into the u variable
//     err := json.NewDecoder(r.Body).Decode(&u)
//     if err != nil {
//         // If there's an error decoding the JSON, return an HTTP 400 error
//         http.Error(w, err.Error(), http.StatusBadRequest)
//         return
//     }

//     // Generate a new ObjectID for the user
//     u.ID = bson.NewObjectId()

//     // Insert the user into the MongoDB database
//     err = uc.session.DB("mongogolang").C("users").Insert(u)
//     if err != nil {
//         // If there's an error inserting the user, return an HTTP 500 error
//         http.Error(w, err.Error(), http.StatusInternalServerError)
//         return
//     }

//     // Marshal the user object into JSON format
//     userJSON, err := json.Marshal(u)
//     if err != nil {
//         // If there's an error marshaling the user object, return an HTTP 500 error
//         http.Error(w, err.Error(), http.StatusInternalServerError)
//         return
//     }

//     // Set the Content-Type header to indicate that JSON data is being returned
//     w.Header().Set("Content-Type", "application/json")

//     // Set the HTTP status code to 201 (Created)
//     w.WriteHeader(http.StatusCreated)

//     // Write the JSON response containing the created user data to the response writer
//     fmt.Fprintf(w, "%s\n", userJSON)
// }

// func (uc UserController) DeleteUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
//     // Get the ID parameter from the request URL
//     id := p.ByName("id")

//     // Convert the ID parameter to bson.ObjectId type
//     oid := bson.ObjectIdHex(id)

//     // Remove the user with the specified ID from the MongoDB database
//     err := uc.session.DB("mongogolang").C("users").RemoveId(oid)
//     if err != nil {
//         // If there's an error removing the user, return an HTTP 404 error
//         http.Error(w, err.Error(), http.StatusNotFound)
//         return
//     }

//     // Set the HTTP status code to 200 (OK)
//     w.WriteHeader(http.StatusOK)

//     // Write a success message to the response writer
//     fmt.Fprintf(w, "Deleted user: %s\n", oid)
// }

package controllers

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/bson"
    "github.com/julienschmidt/httprouter"
    "github.com/monster0freason/Golang-And-MongoDB-REST-API/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserController represents the controller for handling user-related requests
type UserController struct {
    client *mongo.Client
}

// NewUserController creates a new instance of UserController
func NewUserController(c *mongo.Client) *UserController {
    return &UserController{c}
}

// GetUser retrieves a user by ID
func (uc UserController) GetUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    id := p.ByName("id")

    oid, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        w.WriteHeader(http.StatusNotFound)
        return
    }

    collection := uc.client.Database("mongogolang").Collection("users")
    ctx := context.Background()
    filter := bson.M{"_id": oid}

	
    var u models.User
    err = collection.FindOne(ctx, filter).Decode(&u)
    if err != nil {
        w.WriteHeader(http.StatusNotFound)
        return
    }

    uJSON, err := json.Marshal(u)
    if err != nil {
        fmt.Println(err)
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "%s\n", uJSON)
}

func (uc UserController) CreateUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    // Initialize a variable u of type models.User
    var u models.User

    // Decode the JSON data from the request body into the u variable
    err := json.NewDecoder(r.Body).Decode(&u)
    if err != nil {
        // If there's an error decoding the JSON, return an HTTP 400 error
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Set the ID field to a new ObjectID
    u.ID = primitive.NewObjectID()	

    // Insert the user into the MongoDB database
    collection := uc.client.Database("mongogolang").Collection("users")
    ctx := context.Background()

    // Omit setting the ID field, allowing MongoDB to generate a new _id
    _, err = collection.InsertOne(ctx, u)
    if err != nil {
        // If there's an error inserting the user, return an HTTP 500 error
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Marshal the user object into JSON format
    userJSON, err := json.Marshal(u)
    if err != nil {
        // If there's an error marshaling the user object, return an HTTP 500 error
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Set the Content-Type header to indicate that JSON data is being returned
    w.Header().Set("Content-Type", "application/json")

    // Set the HTTP status code to 201 (Created)
    w.WriteHeader(http.StatusCreated)

    // Write the JSON response containing the created user data to the response writer
    fmt.Fprintf(w, "%s\n", userJSON)
}


func (uc UserController) DeleteUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    // Get the ID parameter from the request URL
    id := p.ByName("id")

    oid, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    // Remove the user with the specified ID from the MongoDB database
    collection := uc.client.Database("mongogolang").Collection("users")
    ctx := context.Background()
    filter := bson.M{"_id": oid}
    _, err = collection.DeleteOne(ctx, filter)
    if err != nil {
        // If there's an error removing the user, return an HTTP 404 error
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }

    // Set the HTTP status code to 200 (OK)
    w.WriteHeader(http.StatusOK)

    // Write a success message to the response writer
    fmt.Fprintf(w, "Deleted user: %s\n", oid)
}
