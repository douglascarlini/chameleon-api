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

var DB_HOST = os.Getenv("DB_HOST")
var DB_PORT = os.Getenv("DB_PORT")
var DB_NAME = os.Getenv("DB_NAME")
var DB_USER = os.Getenv("DB_USER")
var DB_PASS = os.Getenv("DB_PASS")

func init() {
	if len(DB_HOST) == 0 {
		DB_HOST = "localhost"
	}
	if len(DB_PORT) == 0 {
		DB_PORT = "27017"
	}
	if len(DB_USER) == 0 {
		DB_USER = "root"
	}
	if len(DB_PASS) == 0 {
		DB_PASS = "root"
	}
	if len(DB_NAME) == 0 {
		DB_NAME = "data"
	}
}

var client *mongo.Client

func main() {
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s:%s", DB_USER, DB_PASS, DB_HOST, DB_PORT))

	var err error
	client, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/{collection}", createHandler).Methods("POST")
	r.HandleFunc("/{collection}/{id}", updateHandler).Methods("PUT")
	r.HandleFunc("/{collection}/{id}", deleteHandler).Methods("DELETE")
	r.HandleFunc("/{collection}/search", searchHandler).Methods("POST")

	http.Handle("/", r)
	log.Print("Running on port 80...")
	log.Fatal(http.ListenAndServe(":80", r))
}
