# Golang And MongoDB REST API
 "Golang CRUD API with MongoDB integration. Implements basic user management functionalities: get, create, delete. Simplified structure for beginners."
## 0️⃣ Overview of the project


## 3️⃣ Library Imports

1. **HTTP Router Package (github.com/julienschmidt/httprouter):**
   
   This package provides a lightweight and high-performance HTTP request router for Go. It's commonly used for building web applications and APIs due to its efficiency and simplicity. The `httprouter` package is known for its fast routing capabilities, making it suitable for high-traffic applications.

   ```bash
   go get github.com/julienschmidt/httprouter
   ```

2. **MongoDB Driver Package (go.mongodb.org/mongo-driver/mongo):**
   
   The `mongo-driver` package is the official MongoDB driver for Go, offering support for the latest MongoDB features and performance improvements. It allows Go applications to interact with MongoDB databases by providing a comprehensive set of functions and types for CRUD operations, querying data, and managing connections. This package enables seamless integration between Go applications and MongoDB databases.

   ```bash
   go get go.mongodb.org/mongo-driver/mongo
   ```

3. **BSON Package (go.mongodb.org/mongo-driver/bson):**
   
   BSON (Binary JSON) is the binary serialization format used by MongoDB to store and transmit data. The `bson` package, located within the `mongo-driver` package, provides functionality to work with BSON data in Go. It allows developers to serialize Go data structures into BSON format and vice versa, facilitating communication between Go applications and MongoDB databases.

   ```bash
   go get go.mongodb.org/mongo-driver/bson
   ```

4. **Primitive BSON Package (go.mongodb.org/mongo-driver/bson/primitive):**
   
   The `primitive` package, also part of the `mongo-driver` package, offers support for primitive BSON data types such as ObjectID, Timestamp, Binary, and Decimal128. It provides methods for creating, parsing, and manipulating these primitive BSON types, which are commonly used in MongoDB documents.

   ```bash
   go get go.mongodb.org/mongo-driver/bson/primitive
   ```

These commands will download and install the specified packages and their dependencies into your Go environment, allowing you to utilize their functionality in your Go applications. With these packages installed, you can proceed to develop your Go application, implementing HTTP routing with `httprouter` and integrating MongoDB functionality using the `mongo-driver` package along with `bson` and `primitive` sub-packages.




## 4️⃣ Project Structure
```
- Project_Directory/
    - main.go
    - controllers/
        - user.go
    - models/
        - user.go
```

Explanation:

1. **Project_Directory:** This is the root directory of your project where all files and folders will reside.

2. **main.go:** This file serves as the entry point for your Go application. It contains the main function and handles routing logic for HTTP requests.

3. **controllers/:** This folder contains controller files responsible for handling HTTP requests and responses. Each controller file corresponds to specific HTTP endpoints.

    - **user.go:** This file contains functions related to user operations such as creating, reading, updating, and deleting users. It handles HTTP requests for user-related operations.

4. **models/:** This folder contains model files representing data structures used in your application. Each model file defines a Go struct representing a data entity.

    - **user.go:** This file defines the structure of the User model, including fields such as ID, Name, Email, etc.

This structure separates concerns by placing HTTP request handling logic in controllers and data structure definitions in models. It provides scalability, readability, and modularity to your project, making it easier to understand and maintain as it grows.


## 5️⃣ `main.go` file 

1. **main.go:**
```go
package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Main function, entry point of the application
func main() {
	// Create a new instance of HTTP router
	r := httprouter.New()

	// Define routes and associate them with corresponding controller functions
	r.GET("/user/:id", uc.GetUser)
	r.POST("/user", uc.CreateUser)
	r.DELETE("/user/:id", uc.DeleteUser)

	// Start the HTTP server on port 9000
	http.ListenAndServe(":9000", r)
}

// getSession establishes a connection to MongoDB and returns a session
func getSession() *mgo.Session {
	// Connect to MongoDB
	s, err := mgo.Dial("localhost:27317")
	if err != nil {
		panic(err)
	}
	return s
}
```

- **Explanation:**
  - This file serves as the entry point of the application.
  - We import necessary packages: `net/http` for handling HTTP requests and responses, and `github.com/julienschmidt/httprouter` for routing.
  - In the `main` function, we create a new instance of `httprouter.Router`.
  - We define three routes using the `GET`, `POST`, and `DELETE` methods:
    - `GET /user/:id`: Retrieve user details by ID.
    - `POST /user`: Create a new user.
    - `DELETE /user/:id`: Delete a user by ID.
  - These routes are associated with corresponding controller functions (`uc.GetUser`, `uc.CreateUser`, `uc.DeleteUser`), which will handle the logic for processing the requests.
  - We start the HTTP server on port `9000` using `http.ListenAndServe`.
  - The `getSession` function establishes a connection to MongoDB using `mgo.Dial`. It returns a session that can be used for interacting with the database.




## 4

2. **models/user.go:**
```go
package models

import "gopkg.in/mgo.v2/bson"

// User represents the user model
type User struct {
	ID     bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name   string        `json:"name" bson:"name"`
	Gender string        `json:"gender" bson:"gender"`
	Age    int           `json:"age" bson:"age"`
}

// Other fields and functions can be added as needed
```

- **Explanation:**
  - This file defines the `User` model representing a user in the application.
  - The `User` struct contains fields such as `ID`, `Name`, `Gender`, and `Age`, which represent properties of a user.
  - Each field is tagged with `json` to specify how it should be serialized to JSON when sending responses, and `bson` to specify how it should be stored in MongoDB.
  - For example, the `ID` field is tagged as `json:"id" bson:"_id,omitempty"`, which means it will be serialized as `id` in JSON and stored as `_id` in MongoDB.

Overall, these files lay the foundation for a Go web application with routing, database interaction, and data modeling capabilities. The `main.go` file handles HTTP routing and server setup, while the `models/user.go` file defines the data structure for representing users in the application.


## 5
Sure, let's write the code for the `controllers/user.go` file and explain it in detail:

```go
package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2/bson"
	"github.com/yourusername/projectname/models"
)

// UserController represents the controller for handling user-related requests
type UserController struct {
	session *mgo.Session
}

// NewUserController creates a new instance of UserController
func NewUserController(s *mgo.Session) *UserController {
	return &UserController{s}
}

// GetUser retrieves a user by ID
func (uc UserController) GetUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")

	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	oid := bson.ObjectIdHex(id)

	u := models.User{}

	if err := uc.session.DB("mongogolang").C("users").FindId(oid).One(&u); err != nil {
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
```

- **Explanation:**
  - We start by defining the package and importing necessary packages, including `encoding/json`, `fmt`, `net/http`, `github.com/julienschmidt/httprouter`, `gopkg.in/mgo.v2/bson`, and `github.com/yourusername/projectname/models`.
  - `UserController` is a struct representing the controller for handling user-related requests. It contains a session pointer to MongoDB.
  - `NewUserController` is a constructor function that creates a new instance of `UserController` with the provided MongoDB session.
  - `GetUser` is a method of `UserController` that handles the retrieval of a user by ID.
    - It takes `http.ResponseWriter`, `*http.Request`, and `httprouter.Params` as parameters.
    - It extracts the ID from the request parameters and checks if it's a valid MongoDB ObjectID. If not, it returns a 404 status code.
    - It converts the ID to a MongoDB `ObjectId`.
    - It initializes an empty `models.User` struct.
    - It queries the MongoDB database for the user with the specified ID.
    - If the user is not found, it returns a 404 status code.
    - If the user is found, it marshals the user data into JSON format.
    - If marshaling fails, it returns a 500 status code.
    - It sets the response content type to JSON, writes a 200 status code, and sends the JSON data in the response body.

This controller file defines the logic for handling user-related HTTP requests. The `GetUser` method retrieves a user by ID from the database and returns the user data in JSON format. The controller is responsible for interacting with the database and formatting the response accordingly.

Here's a breakdown of the `CreateUser` and `UserController` functions along with explanations:

### `CreateUser` Function:

```go
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

    // Generate a new ObjectID for the user
    u.ID = bson.NewObjectId()

    // Insert the user into the MongoDB database
    err = uc.session.DB("mongogolang").C("users").Insert(u)
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
```

#### Explanation:
1. The `CreateUser` function takes three parameters: `w` (http.ResponseWriter), `r` (http.Request), and `_` (httprouter.Params). The `_` parameter is used to indicate that we're not using the httprouter.Params parameter in this function.
  
2. A variable `u` of type `models.User` is initialized to store the user data.

3. The JSON data from the request body (`r.Body`) is decoded into the `u` variable.

4. A new ObjectID is generated for the user.

5. The user data is inserted into the MongoDB database.

6. The inserted user object is marshaled into JSON format.

7. The Content-Type header is set to indicate that JSON data is being returned.

8. The HTTP status code is set to 201 (Created).

9. The JSON response containing the created user data is written to the response writer.

### `UserController` Function:

```go
func (uc UserController) UserController(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    // Get the ID parameter from the request URL
    id := p.ByName("id")

    // Convert the ID parameter to bson.ObjectId type
    oid := bson.ObjectIdHex(id)

    // Remove the user with the specified ID from the MongoDB database
    err := uc.session.DB("mongogolang").C("users").RemoveId(oid)
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
```

#### Explanation:
1. The `UserController` function takes three parameters: `w` (http.ResponseWriter), `r` (http.Request), and `p` (httprouter.Params). The `p` parameter contains the route parameters, including the user ID.

2. The ID parameter is retrieved from the route parameters using `ByName`.

3. The ID parameter is converted to a `bson.ObjectId` type.

4. The user with the specified ID is removed from the MongoDB database.

5. If there's an error removing the user, an HTTP 404 error is returned.

6. If the user is successfully deleted, an HTTP 200 (OK) status code is returned along with a success message indicating the deleted user's ID.

These functions together handle the creation and deletion of users in a MongoDB database using Go.