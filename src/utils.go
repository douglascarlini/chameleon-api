package main

import (
	"encoding/json"
	"net/http"
)

func sendData(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	obj, _ := json.Marshal(data)
	w.WriteHeader(http.StatusOK)
	w.Write(obj)
}

func sendBadRequest(w http.ResponseWriter, message string) {
	http.Error(w, message, http.StatusBadRequest)
}

func sendUnauthorized(w http.ResponseWriter, message string) {
	http.Error(w, message, http.StatusUnauthorized)
}

func sendNotFound(w http.ResponseWriter, message string) {
	http.Error(w, message, http.StatusNotFound)
}

func sendConflict(w http.ResponseWriter, message string) {
	http.Error(w, message, http.StatusConflict)
}

func sendError(w http.ResponseWriter, message string) {
	http.Error(w, message, http.StatusInternalServerError)
}
