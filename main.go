package main

// import (
// 	"net/http"
// 	"gopkg.in/mgo.v2"
// 	"github.com/julienschmidt/httprouter"
// 	"github.com/monster0freason/Golang-And-MongoDB-REST-API/controllers"
	
// )

// // Main function, entry point of the application
// func main() {
// 	// Create a new instance of HTTP router
// 	r := httprouter.New()

// 	uc := controllers.NewUserController(getSession())

// 	// Define routes and associate them with corresponding controller functions
// 	r.GET("/user/:id", uc.GetUser)
// 	r.POST("/user", uc.CreateUser)
// 	r.DELETE("/user/:id", uc.DeleteUser)

// 	// Start the HTTP server on port 9000
// 	http.ListenAndServe("localhost:9000", r)
// }

// // getSession establishes a connection to MongoDB and returns a session
// func getSession() *mgo.Session {
// 	// Connect to MongoDB
// 	s, err := mgo.Dial("localhost:27017")
// 	if err != nil {
// 		panic(err)
// 	}
// 	return s
// }





import (
	"context"
	"github.com/monster0freason/Golang-And-MongoDB-REST-API/controllers"
	"net/http"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	
	
)

// Main function, entry point of the application
func main() {
	// Create a new instance of HTTP router
	r := httprouter.New()

	uc := controllers.NewUserController(getMongoClient())

	// Define routes and associate them with corresponding controller functions
	r.GET("/user/:id", uc.GetUser)
	r.POST("/user", uc.CreateUser)
	r.DELETE("/user/:id", uc.DeleteUser)

	// Start the HTTP server on port 9000
	http.ListenAndServe("localhost:9000", r)
}

// getMongoClient establishes a connection to MongoDB and returns a MongoDB client
func getMongoClient() *mongo.Client {
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		panic(err)
	}

	// Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		panic(err)
	}

	return client
}
