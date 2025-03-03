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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	collection := db.Collection(collectionName)
	res, err := collection.DeleteOne(context.Background(), bson.M{"_id": oid})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]any{"rows": res.DeletedCount}
	respondWithJSON(w, data, http.StatusOK)
}
