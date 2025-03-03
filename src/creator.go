package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WriteData struct {
	Name     string         `json:"name"`
	Data     map[string]any `json:"data"`
	ClientID string         `json:"client_id"`
}

func Creator() {

	err = rmq.NewQueue("write")
	if err != nil {
		log.Fatal(err)
	}

	rmq.Pipe["write"].Consume(func(b []byte, id string) {

		var data WriteData

		fmt.Printf("TICKET: %s\n", id)

		err = json.Unmarshal(b, &data)
		if err != nil {
			log.Fatal(err)
			return
		}

		for key, value := range data.Data {
			if strings.HasSuffix(key, "_id") {
				value, isString := value.(string)
				if isString {
					oid, err := primitive.ObjectIDFromHex(value)
					if err == nil {
						data.Data[key] = oid
					}
					dt, err := time.Parse("2006-01-02 15:04:05", value)
					if err == nil {
						data.Data[key] = dt
					}
				}
			}
		}

		collection := db.Collection(data.Name)
		res, err := collection.InsertOne(context.Background(), data.Data)
		if err != nil {
			log.Fatal(err)
			return
		}

		if len(data.ClientID) > 0 {
			msg := map[string]any{"ticket": id, "result": res.InsertedID}
			err = sendMessageToClient(data.ClientID, msg)
			if err != nil {
				fmt.Printf("ws error: %s", err.Error())
			}
		}

	})

}
