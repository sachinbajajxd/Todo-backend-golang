package controllers

import (
	"context"
	"crud/models"
	"encoding/json"
	"fmt"
	"strconv"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("backend-crud")

type Claims struct {
    Username string `json:"username"`
    jwt.StandardClaims
}

func RegisterHandler(client *mongo.Client, w http.ResponseWriter, r *http.Request) {
    var user models.User
    err := json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }

    // Checking if user already exists
    collection := client.Database("test").Collection("users")
    var existingUser models.User
    err = collection.FindOne(context.TODO(), bson.M{"username": user.Username}).Decode(&existingUser)
    if err == nil {
        http.Error(w, "Username already exists", http.StatusConflict)
        return
    } else if err != mongo.ErrNoDocuments {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Hashing the password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    user.Password = string(hashedPassword)

    // Generating unique ID
    user.ID = primitive.NewObjectID()

    _, err = collection.InsertOne(context.TODO(), user)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}

func LoginHandler(client *mongo.Client, w http.ResponseWriter, r *http.Request) {
    var user models.User
    err := json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }

    // Checking if user exists
    collection := client.Database("test").Collection("users")
    var foundUser models.User
    err = collection.FindOne(context.TODO(), bson.M{"username": user.Username}).Decode(&foundUser)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            http.Error(w, "Invalid username or password", http.StatusUnauthorized)
            return
        }
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Comparing the password with the hash
    err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(user.Password))
    if err != nil {
        http.Error(w, "Invalid username or password", http.StatusUnauthorized)
        return
    }

    // Generating JWT token
	token, err := GenerateJWT(user.Username)
    if err != nil {
        http.Error(w, "Error generating token", http.StatusInternalServerError)
        return
    }

	type LoginResponse struct {
		Token string      `json:"token"`
		User  models.User `json:"user"`
	}

	response := LoginResponse{
		Token: token,
		User: foundUser,
	}

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(response)
}

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
    client := r.Context().Value("mongoClient").(*mongo.Client)

    var post models.Post
    err := json.NewDecoder(r.Body).Decode(&post)
    if err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    post.ID = primitive.NewObjectID()
    post.CreatedAt = time.Now()
    post.UpdatedAt = time.Now()

    // Create a new post
    collection := client.Database("test").Collection("posts")
    _, err = collection.InsertOne(context.TODO(), post)
    if err != nil {
        http.Error(w, "Failed to create post", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(post)
}

func GetAllPostHandler(w http.ResponseWriter, r *http.Request) {
    client := r.Context().Value("mongoClient").(*mongo.Client)

    // Retrieving query parameters for sorting and pagination
    queryParams := r.URL.Query()
    userID := queryParams.Get("userId")
    sortBy := queryParams.Get("sortBy")
    sortOrder := queryParams.Get("sortOrder")
    page := queryParams.Get("page")
    limit := queryParams.Get("limit")

    // Default values for pagination
    defaultPage := 1
    defaultLimit := 10

    // Converting page and limit to integers
    pageInt, err := strconv.Atoi(page)
    if err != nil || pageInt <= 0 {
        pageInt = defaultPage
    }

    limitInt, err := strconv.Atoi(limit)
    if err != nil || limitInt <= 0 {
        limitInt = defaultLimit
    }

    sortOrderInt := 1 // ascending order
    if sortOrder == "desc" {
        sortOrderInt = -1
    }

    findOptions := options.Find()
    findOptions.SetSkip(int64((pageInt - 1) * limitInt))
    findOptions.SetLimit(int64(limitInt))
    findOptions.SetSort(bson.D{{Key: sortBy, Value: sortOrderInt}})

    userObjID, e := primitive.ObjectIDFromHex(userID)
    if e != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    filter := bson.M{"user_id": userObjID}

    var posts []models.Post

    cursor, err := client.Database("test").Collection("posts").Find(context.TODO(), filter, findOptions)
    if err != nil {
        http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
        return
    }
    defer cursor.Close(context.TODO())

    for cursor.Next(context.TODO()) {
        var post models.Post
        if err := cursor.Decode(&post); err != nil {
            http.Error(w, "Failed to decode post", http.StatusInternalServerError)
            return
        }
        posts = append(posts, post)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(posts)
}

func UpdatePostHandler(w http.ResponseWriter, r *http.Request) {
    client := r.Context().Value("mongoClient").(*mongo.Client)
    vars := mux.Vars(r)
    postID := vars["id"]

    postObjID, e := primitive.ObjectIDFromHex(postID)
    if e != nil {
        http.Error(w, "Invalid post ID", http.StatusBadRequest)
        return
    }
    var updateRequest models.UpdatePostRequest
    err := json.NewDecoder(r.Body).Decode(&updateRequest)
    if err != nil {
        fmt.Println(err)
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    // Checking if post exists and user has access
    userID := updateRequest.UserID
    var post models.Post
    err = client.Database("test").Collection("posts").FindOne(context.TODO(), bson.M{"_id": postObjID, "user_id": userID}).Decode(&post)
    if err != nil {
        fmt.Println(err)
        http.Error(w, "Post not found or access denied", http.StatusNotFound)
        return
    }

    updateFields := bson.M{}
    if updateRequest.Title != nil {
        updateFields["title"] = *updateRequest.Title
    }
    if updateRequest.Description != nil {
        updateFields["description"] = *updateRequest.Description
    }
    if updateRequest.Status != nil {
        updateFields["status"] = *updateRequest.Status
    }
    updateFields["updated"] = time.Now()

    _, err = client.Database("test").Collection("posts").UpdateOne(context.TODO(), bson.M{"_id": postObjID}, bson.M{"$set": updateFields})
    if err != nil {
        http.Error(w, "Failed to update post", http.StatusInternalServerError)
        return
    }

    updatedPost := models.Post{}
    err = client.Database("test").Collection("posts").FindOne(context.TODO(), bson.M{"_id": postObjID}).Decode(&updatedPost)
    if err != nil {
        http.Error(w, "Failed to fetch updated post", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(updatedPost)
}


func DeletePostHandler(w http.ResponseWriter, r *http.Request) {
    client := r.Context().Value("mongoClient").(*mongo.Client)
    vars := mux.Vars(r)
    postID := vars["id"]

    postObjID, e := primitive.ObjectIDFromHex(postID)
    if e != nil {
        http.Error(w, "Invalid post ID", http.StatusBadRequest)
        return
    }
    var deleteRequest models.DeletePostRequest
    err := json.NewDecoder(r.Body).Decode(&deleteRequest)
    if err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    userID := deleteRequest.UserID
    if userID == primitive.NilObjectID {
        http.Error(w, "User ID is required", http.StatusBadRequest)
        return
    }
    
    // Checking if post exists
    var post models.Post
    err = client.Database("test").Collection("posts").FindOne(context.TODO(), bson.M{"_id": postObjID, "user_id": userID}).Decode(&post)
    if err != nil {
        http.Error(w, "Post not found or access denied", http.StatusNotFound)
        return
    }

    // Deleting the post
    _, err = client.Database("test").Collection("posts").DeleteOne(context.TODO(), bson.M{"_id": postObjID})
    if err != nil {
        http.Error(w, "Failed to delete post", http.StatusInternalServerError)
        return
    }

    successResponse := map[string]bool{"success": true}
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(successResponse)
}


// Middleware to generate JWT token
func GenerateJWT(username string) (string, error) {
    expirationTime := time.Now().Add(24 * time.Hour)
    claims := &Claims{
        Username: username,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: expirationTime.Unix(),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(jwtKey)
    if err != nil {
        return "", err
    }

    return tokenString, nil
}


func Authenticate(next http.Handler, client *mongo.Client) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        tokenStr := r.Header.Get("Authorization")
        if tokenStr == "" {
            http.Error(w, "Missing token", http.StatusUnauthorized)
            return
        }

        // Split the tokenStr to remove "Bearer" part
        parts := strings.Split(tokenStr, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            http.Error(w, "Invalid token format", http.StatusUnauthorized)
            return
        }
        tokenStr = parts[1]

        claims := &Claims{}
        token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
            return jwtKey, nil
        })
        if err != nil || !token.Valid {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }

        ctx := context.WithValue(r.Context(), "mongoClient", client)
        r = r.WithContext(ctx)
        next.ServeHTTP(w, r)
    })
}