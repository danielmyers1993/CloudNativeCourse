package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
)

var mutex sync.RWMutex

func main() {
	db := database{"shoes": 50, "socks": 5}
	mux := http.NewServeMux()
	mux.HandleFunc("/list", db.list)
	mux.HandleFunc("/price", db.price)
	mux.HandleFunc("/create", db.create)
	mux.HandleFunc("/update", db.update)
	mux.HandleFunc("/delete", db.delete)
	log.Fatal(http.ListenAndServe(":8000", mux))
}

type dollars float32

func (d dollars) String() string { return fmt.Sprintf("$%.2f", d) }

type database map[string]dollars

func (db database) list(w http.ResponseWriter, req *http.Request) {
	mutex.RLock()
	defer mutex.RUnlock()

	for item, price := range db {
		fmt.Fprintf(w, "%s: %s\n", item, price)
	}
}

func (db database) price(w http.ResponseWriter, req *http.Request) {
	mutex.RLock()
	defer mutex.RUnlock()

	item := req.URL.Query().Get("item")
	if price, ok := db[item]; ok {
		fmt.Fprintf(w, "%s\n", price)
	} else {
		w.WriteHeader(http.StatusNotFound) // 404
		fmt.Fprintf(w, "no such item: %q\n", item)
	}
}

func (db database) create(w http.ResponseWriter, req *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	item := req.URL.Query().Get("item")
	priceStr := req.URL.Query().Get("price")

	price, err := strconv.ParseFloat(priceStr, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest) // 400
		fmt.Fprintf(w, "invalid price format: %s\n", priceStr)
		return
	}

	db[item] = dollars(price)
	fmt.Fprintf(w, "item %s created successfully\n", item)
}

func (db database) update(w http.ResponseWriter, req *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	item := req.URL.Query().Get("item")
	newPriceStr := req.URL.Query().Get("price")

	newPrice, err := strconv.ParseFloat(newPriceStr, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest) // 400
		fmt.Fprintf(w, "invalid price format: %s\n", newPriceStr)
		return
	}

	if _, ok := db[item]; !ok {
		w.WriteHeader(http.StatusNotFound) // 404
		fmt.Fprintf(w, "no such item: %q\n", item)
		return
	}

	db[item] = dollars(newPrice)
	fmt.Fprintf(w, "item %s updated successfully\n", item)
}

func (db database) delete(w http.ResponseWriter, req *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	item := req.URL.Query().Get("item")

	if _, ok := db[item]; !ok {
		w.WriteHeader(http.StatusNotFound) // 404
		fmt.Fprintf(w, "no such item: %q\n", item)
		return
	}

	delete(db, item)
	fmt.Fprintf(w, "item %s deleted successfully\n", item)
}
