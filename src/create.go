package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func createHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
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

	for key, value := range payload {
		if strings.HasSuffix(key, "_id") {
			value, isString := value.(string)
			if isString {
				oid, err := primitive.ObjectIDFromHex(value)
				if err == nil {
					payload[key] = oid
				}
				dt, err := time.Parse("2006-01-02 15:04:05", value)
				if err == nil {
					payload[key] = dt
				}
			}
		}
	}

	payload["created"] = time.Now()

	collection := client.Database(DB_NAME).Collection(collectionName)
	res, err := collection.InsertOne(context.Background(), payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{"id": res.InsertedID}
	respondWithJSON(w, data, http.StatusCreated)
}
