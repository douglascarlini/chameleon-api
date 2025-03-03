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

var client *mongo.Client

func main() {
	dsn := fmt.Sprintf("mongodb://%s:%s@%s:%s", DB_USER, DB_PASS, DB_HOST, DB_PORT)
	clientOptions := options.Client().ApplyURI(dsn)

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
