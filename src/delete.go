package main

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	collectionName := params["collection"]
	id := params["id"]

	collection := client.Database(DB_NAME).Collection(collectionName)
	res, err := collection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{"rows": res.DeletedCount}
	respondWithJSON(w, data, http.StatusOK)
}
