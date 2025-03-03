package main

import (
	"context"
	"fmt"
	"os"

	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var JWT_CODE = os.Getenv("JWT_CODE")
var MDB_NAME = os.Getenv("MDB_NAME")
var MDB_USER = os.Getenv("MDB_USER")
var MDB_PASS = os.Getenv("MDB_PASS")
var db *mongo.Database

var client *mongo.Client

func main() {
	dsn := fmt.Sprintf("mongodb://%s:%s@mdb:27017", MDB_USER, MDB_PASS)
	clientOptions := options.Client().ApplyURI(dsn)

	var err error
	client, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	db = client.Database(MDB_NAME)

	r := mux.NewRouter()

	auth := r.PathPrefix("/auth").Subrouter()
	api := r.PathPrefix("/api").Subrouter()
	api.Use(jwtAuthMiddleware)

	auth.HandleFunc("/signup", signUp).Methods("POST")
	auth.HandleFunc("/signin", signIn).Methods("POST")

	api.HandleFunc("/{collection}", createHandler).Methods("POST")
	api.HandleFunc("/{collection}/{id}", updateHandler).Methods("PUT")
	api.HandleFunc("/{collection}/{id}", deleteHandler).Methods("DELETE")
	api.HandleFunc("/{collection}/search", searchHandler).Methods("POST")

	http.Handle("/", r)
	log.Print("Running on port 80...")
	log.Fatal(http.ListenAndServe(":80", r))
}
