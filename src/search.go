package main

import (
	"context"
	"encoding/json"
	"io"
	"strings"

	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func searchHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	collectionName := params["collection"]

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var searchOptions map[string]interface{}
	err = json.Unmarshal(body, &searchOptions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	paging, _ := searchOptions["paging"].(map[string]interface{})
	filter, _ := searchOptions["filter"].(map[string]interface{})
	sort, _ := searchOptions["sort"].(map[string]interface{})

	filterQuery := bson.M{}

	if filter != nil {
		filterQuery = filter
	}

	collection := client.Database(DB_NAME).Collection(collectionName)

	var sortOptions bson.D
	for key, value := range sort {
		sortOptions = append(sortOptions, bson.E{Key: key, Value: value})
	}

	options := options.Find()
	if len(sortOptions) > 0 {
		options.SetSort(sortOptions)
	}

	if paging != nil {
		page, _ := paging["page"].(float64)
		limit, _ := paging["limit"].(float64)

		if page > 0 && limit > 0 {
			options.SetSkip(int64((page - 1) * limit))
			options.SetLimit(int64(limit))
		}
	}

	cursor, err := collection.Find(context.Background(), filterQuery, options)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	var results []map[string]interface{}
	for cursor.Next(context.Background()) {
		var result map[string]interface{}
		err := cursor.Decode(&result)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var joinOptions map[string]interface{}
		joinOptions, hasJoins := searchOptions["join"].(map[string]interface{})

		if hasJoins {
			for joinCollection, joinField := range joinOptions {

				joinFilter := bson.M{joinField.(string): result["_id"]}
				joinCursor, err := client.Database(DB_NAME).Collection(joinCollection).Find(context.Background(), joinFilter)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				defer joinCursor.Close(context.Background())

				var joinResults []map[string]interface{}
				for joinCursor.Next(context.Background()) {
					var joinResult map[string]interface{}
					err := joinCursor.Decode(&joinResult)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
					joinResults = append(joinResults, joinResult)
				}

				result[joinCollection] = joinResults
			}
		}

		loadOptions, hasLoad := searchOptions["load"].(map[string]interface{})

		if hasLoad {
			for targetCollection, fieldMapping := range loadOptions {
				fieldOnCollection := strings.Split(fieldMapping.(string), ":")
				result[fieldOnCollection[0]] = loadEntryOnCollection(targetCollection, result[fieldOnCollection[1]])
			}
		}

		results = append(results, result)
	}

	respondWithJSON(w, results, http.StatusOK)
}

func loadEntryOnCollection(collection string, value interface{}) interface{} {
	targetColl := client.Database(DB_NAME).Collection(collection)

	cursor, err := targetColl.Find(context.Background(), bson.M{"_id": value})
	if err != nil {
		return err
	}
	defer cursor.Close(context.Background())

	if cursor.Next(context.Background()) {
		var record map[string]interface{}
		if err := cursor.Decode(&record); err != nil {
			return err
		}
		return record
	}

	return nil
}
