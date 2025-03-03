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
		sendBadRequest(w, err.Error())
		return
	}

	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		sendBadRequest(w, err.Error())
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

	collection := db.Collection(collectionName)
	res, err := collection.InsertOne(context.Background(), payload)
	if err != nil {
		sendError(w, err.Error())
		return
	}

	sendData(w, map[string]any{"_id": res.InsertedID})
}
