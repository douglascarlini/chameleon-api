package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func signIn(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		sendBadRequest(w, err.Error())
		return
	}

	var user User
	if err := json.Unmarshal(body, &user); err != nil {
		sendBadRequest(w, err.Error())
		return
	}

	var load map[string]any
	collection := db.Collection("users")
	filter := map[string]any{"username": user.Username}

	err = collection.FindOne(context.Background(), filter).Decode(&load)

	if err != nil {
		if err != mongo.ErrNoDocuments {
			sendUnauthorized(w, "")
			return
		}
		sendError(w, err.Error())
		return
	}

	if checkPassword(user.Password, load["password"].(string)) {

		_id := load["_id"].(primitive.ObjectID).Hex()

		token, err := generateToken(_id)
		if err != nil {
			sendError(w, err.Error())
			return
		}

		sendData(w, map[string]any{
			"token": token,
			"user": map[string]any{
				"_id":      _id,
				"name":     load["name"].(string),
				"phone":    load["phone"].(string),
				"email":    load["email"].(string),
				"username": load["username"].(string),
			},
		})
		return

	}

	sendUnauthorized(w, "")

}
