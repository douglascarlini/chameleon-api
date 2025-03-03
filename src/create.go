package main

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func createHandler(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	collectionName := params["collection"]

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	payload["created"] = time.Now()
	clientID := r.Header.Get("Client-ID")
	msg := WriteData{Name: collectionName, Data: payload, ClientID: clientID}

	id, err := rmq.Pipe["write"].Publish(msg)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]string{"ticket": id}
	respondWithJSON(w, data, http.StatusCreated)

}
