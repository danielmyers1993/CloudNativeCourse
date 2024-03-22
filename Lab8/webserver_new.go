package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Connect(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	db := client.Database("mydatabase").Collection("products")

	mux := http.NewServeMux()
	mux.HandleFunc("/list", listHandler(db))
	mux.HandleFunc("/price", priceHandler(db))

	// Additional Handlers
	mux.HandleFunc("/create", createHandler(db))
	mux.HandleFunc("/read", readHandler(db))
	mux.HandleFunc("/update", updateHandler(db))
	mux.HandleFunc("/delete", deleteHandler(db))

	log.Fatal(http.ListenAndServe(":8000", mux))
}

func listHandler(db *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		cur, err := db.Find(context.Background(), nil)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "error: %v\n", err)
			return
		}
		defer cur.Close(context.Background())

		for cur.Next(context.Background()) {
			var result map[string]interface{}
			err := cur.Decode(&result)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "error: %v\n", err)
				return
			}
			fmt.Fprintf(w, "%s: $%.2f\n", result["item"], result["price"])
		}
	}
}

func priceHandler(db *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		item := req.URL.Query().Get("item")

		var result map[string]interface{}
		err := db.FindOne(context.Background(), map[string]string{"item": item}).Decode(&result)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "no such item: %q\n", item)
			return
		}

		fmt.Fprintf(w, "$%.2f\n", result["price"])
	}
}

func createHandler(db *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		item := req.URL.Query().Get("item")
		newPrice := req.URL.Query().Get("price")

		price, err := strconv.ParseFloat(newPrice, 32)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "invalid price: %q\n", newPrice)
			return
		}

		_, err = db.InsertOne(context.Background(), map[string]interface{}{"item": item, "price": price})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "error: %v\n", err)
			return
		}

		fmt.Fprintf(w, "create item: %s, price: $%.2f\n", item, price)
	}
}

func readHandler(db *mongo.Collection) http.HandlerFunc {
	return listHandler(db)
}

func updateHandler(db *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		item := req.URL.Query().Get("item")
		newPrice := req.URL.Query().Get("price")

		price, err := strconv.ParseFloat(newPrice, 32)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "invalid price: %q\n", newPrice)
			return
		}

		_, err = db.UpdateOne(context.Background(), map[string]string{"item": item}, map[string]interface{}{"$set": map[string]interface{}{"price": price}})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "error: %v\n", err)
			return
		}

		fmt.Fprintf(w, "update item: %s, price: $%.2f\n", item, price)
	}
}

func deleteHandler(db *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		item := req.URL.Query().Get("item")

		_, err := db.DeleteOne(context.Background(), map[string]string{"item": item})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "error: %v\n", err)
			return
		}

		fmt.Fprintf(w, "delete item: %s\n", item)
	}
}

