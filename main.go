package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
	"crud/controllers"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/mongo/readpref"
    "github.com/gorilla/mux"
)

var client *mongo.Client

func main() {
    var err error
    client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb+srv://sachinbajaj0001:DPLr5nGsyRppckCb@cluster0.kr326ac.mongodb.net/"))
    if err != nil {
        log.Fatalf("Error while connecting to MongoDB: %v", err)
    }
    if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
        log.Fatalf("Error while pinging MongoDB: %v", err)
    }
    fmt.Println("Connected to MongoDB")

    //API routes
    router := mux.NewRouter()

    // Register and Login
    router.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
        controllers.RegisterHandler(client, w, r)
    }).Methods("POST")
    router.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
        controllers.LoginHandler(client, w, r)
    }).Methods("POST")

    // Create a post
    router.HandleFunc("/posts", func(w http.ResponseWriter, r *http.Request) {
        controllers.Authenticate(http.HandlerFunc(controllers.CreatePostHandler), client).ServeHTTP(w, r)
    }).Methods("POST")

    // Get all posts
    router.HandleFunc("/posts", func(w http.ResponseWriter, r *http.Request) {
        controllers.Authenticate(http.HandlerFunc(controllers.GetAllPostHandler), client).ServeHTTP(w, r)
    }).Methods("GET")

    // Update a post
    router.HandleFunc("/posts/{id}", func(w http.ResponseWriter, r *http.Request) {
        controllers.Authenticate(http.HandlerFunc(controllers.UpdatePostHandler), client).ServeHTTP(w, r)
    }).Methods("PUT")

    // Delete a post
    router.HandleFunc("/posts/{id}", func(w http.ResponseWriter, r *http.Request) {
        controllers.Authenticate(http.HandlerFunc(controllers.DeletePostHandler), client).ServeHTTP(w, r)
    }).Methods("DELETE")

    log.Println("Server is starting on port 4000...")
    log.Fatal(http.ListenAndServe(":4000", router))
}
