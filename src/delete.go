package main

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	collectionName := params["collection"]
	id := params["id"]

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		sendError(w, err.Error())
		return
	}

	collection := db.Collection(collectionName)
	res, err := collection.DeleteOne(context.Background(), bson.M{"_id": oid})
	if err != nil {
		sendError(w, err.Error())
		return
	}

	sendData(w, map[string]any{"rows": res.DeletedCount})
}
