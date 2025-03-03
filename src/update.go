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
		sendBadRequest(w, err.Error())
		return
	}

	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		sendBadRequest(w, err.Error())
		return
	}

	if len(payload) == 0 {
		sendBadRequest(w, "")
		return
	}

	collection := db.Collection(collectionName)

	oid, err := primitive.ObjectIDFromHex(documentID)
	if err != nil {
		sendBadRequest(w, "Invalid ID")
		return
	}

	filter := bson.M{"_id": oid}
	update := bson.M{"$set": payload}

	res, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		sendError(w, err.Error())
		return
	}

	if res.MatchedCount == 0 {
		sendNotFound(w, "")
		return
	}

	sendData(w, map[string]any{"rows": res.MatchedCount, "affected": res.ModifiedCount})
}
