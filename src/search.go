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

	var searchOptions map[string]any
	err = json.Unmarshal(body, &searchOptions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	paging, _ := searchOptions["paging"].(map[string]any)
	filter, _ := searchOptions["filter"].(map[string]any)
	sort, _ := searchOptions["sort"].(map[string]any)

	filterQuery := bson.M{}

	if filter != nil {
		filterQuery = filter
	}

	collection := db.Collection(collectionName)

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

	var results []map[string]any
	for cursor.Next(context.Background()) {
		var result map[string]any
		err := cursor.Decode(&result)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var joinOptions map[string]any
		joinOptions, hasJoins := searchOptions["join"].(map[string]any)

		if hasJoins {
			for joinCollection, joinField := range joinOptions {

				joinFilter := bson.M{joinField.(string): result["_id"]}
				joinCursor, err := db.Collection(joinCollection).Find(context.Background(), joinFilter)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				defer joinCursor.Close(context.Background())

				var joinResults []map[string]any
				for joinCursor.Next(context.Background()) {
					var joinResult map[string]any
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

		loadOptions, hasLoad := searchOptions["load"].(map[string]any)

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

func loadEntryOnCollection(collection string, value any) any {
	targetColl := db.Collection(collection)

	cursor, err := targetColl.Find(context.Background(), bson.M{"_id": value})
	if err != nil {
		return err
	}
	defer cursor.Close(context.Background())

	if cursor.Next(context.Background()) {
		var record map[string]any
		if err := cursor.Decode(&record); err != nil {
			return err
		}
		return record
	}

	return nil
}
