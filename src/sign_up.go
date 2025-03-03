package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}

func signUp(w http.ResponseWriter, r *http.Request) {
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

	user.CreatedAt = time.Now()

	user.Password, err = hashPassword(user.Password)
	if err != nil {
		sendError(w, err.Error())
		return
	}

	collection := db.Collection("users")

	var load User
	filter := map[string]string{"username": user.Username}
	err = collection.FindOne(context.Background(), filter).Decode(&load)
	if err == nil {
		sendConflict(w, "Username already taken")
		return
	} else if err != mongo.ErrNoDocuments {
		sendError(w, err.Error())
		return
	}

	res, err := collection.InsertOne(context.Background(), &user)
	if err != nil {
		sendError(w, err.Error())
		return
	}

	sendData(w, map[string]any{"_id": res.InsertedID})

}
