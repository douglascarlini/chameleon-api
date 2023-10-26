package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func updateHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	documentID := params["id"]
	collectionName := params["collection"]

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(body, &payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(payload) == 0 {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	collection := client.Database(DB_NAME).Collection(collectionName)

	oid, err := primitive.ObjectIDFromHex(documentID)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	filter := bson.M{"_id": oid}
	update := bson.M{"$set": payload}

	res, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if res.MatchedCount == 0 {
		http.Error(w, "", http.StatusNotFound)
		return
	}

	response := map[string]interface{}{"rows": res.MatchedCount, "affected": res.ModifiedCount}
	respondWithJSON(w, response, http.StatusOK)
}
