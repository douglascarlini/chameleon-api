package main

import (
	"context"
	"fmt"
	"imovelis/queue"
	"sync"

	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	err    error
	rmq    *queue.RMQ
	db     *mongo.Database
	client *mongo.Client

	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	connections = make(map[string]*websocket.Conn)
	mutex       sync.Mutex
)

func main() {

	dsn := fmt.Sprintf("mongodb://%s:%s@mdb:27017", MDB_USER, MDB_PASS)
	clientOptions := options.Client().ApplyURI(dsn)

	client, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	db = client.Database(MDB_NAME)

	dsn = fmt.Sprintf("amqp://%s:%s@rmq:5672", RMQ_USER, RMQ_PASS)
	rmq = queue.Connect(dsn)

	Creator()

	r := mux.NewRouter()
	r.HandleFunc("/ws", wsHandler)
	r.HandleFunc("/{collection}", createHandler).Methods("POST")
	r.HandleFunc("/{collection}/{id}", updateHandler).Methods("PUT")
	r.HandleFunc("/{collection}/{id}", deleteHandler).Methods("DELETE")
	r.HandleFunc("/{collection}/search", searchHandler).Methods("POST")

	http.Handle("/", r)
	log.Print("Running on port 80...")
	log.Fatal(http.ListenAndServe(":80", r))

}
