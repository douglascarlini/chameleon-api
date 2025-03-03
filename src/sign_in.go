package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

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

	var load User
	filter := map[string]any{"username": user.Username}
	collection := db.Collection("users")

	err = collection.FindOne(context.Background(), filter).Decode(&load)
	if err != nil && err != mongo.ErrNoDocuments {
		sendError(w, err.Error())
		return
	}

	if checkPassword(user.Password, load.Password) {

		token, err := generateToken(user.ID)
		if err != nil {
			sendError(w, err.Error())
			return
		}

		sendData(w, map[string]string{"token": token})
		return

	}

	sendUnauthorized(w, "")

}
