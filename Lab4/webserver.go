package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
)

func main() {
	db := database{data: map[string]dollars{"shoes": 50, "socks": 5}}
	mux := http.NewServeMux()
	mux.HandleFunc("/list", db.list)
	mux.HandleFunc("/price", db.price)

	// Additional Handlers
	mux.HandleFunc("/create", db.create)
	mux.HandleFunc("/read", db.read)
	mux.HandleFunc("/update", db.update)
	mux.HandleFunc("/delete", db.delete)

	log.Fatal(http.ListenAndServe(":8000", mux))
}

type dollars float32

func (d dollars) String() string { return fmt.Sprintf("$%.2f", d) }

// Changed Database to Struct for Lock/Unlock
type database struct {
	data map[string]dollars
	mu   sync.RWMutex
}

func (db *database) list(w http.ResponseWriter, req *http.Request) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	for item, price := range db.data {
		fmt.Fprintf(w, "%s: %s\n", item, price)
	}
}

func (db *database) price(w http.ResponseWriter, req *http.Request) {

	item := req.URL.Query().Get("item")

	db.mu.RLock()
	price, ok := db.data[item]
	db.mu.RUnlock()

	if !ok {
		w.WriteHeader(http.StatusNotFound) // 404
		fmt.Fprintf(w, "no such item: %q\n", item)
	}

	fmt.Fprintf(w, "%s\n", price)
}

func (db *database) create(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	newPrice := req.URL.Query().Get("price")

	price, err := strconv.ParseFloat(newPrice, 32)
	// Handle Parsing Failure
	if err != nil {
		w.WriteHeader(http.StatusBadRequest) // 400
		fmt.Fprintf(w, "invalid price: %q\n", newPrice)
		return
	}

	// Lock for Reading
	db.mu.RLock()
	_, ok := db.data[item]
	db.mu.RUnlock()

	if ok {
		w.WriteHeader(http.StatusBadRequest) // 400
		fmt.Fprintf(w, "item already exists: %s\n", item)
		return
	}

	db.mu.Lock()
	(db.data)[item] = dollars(price)
	db.mu.Unlock()

	fmt.Fprintf(w, "create item: %s, price: %s\n", item, (db.data)[item])
}

func (db *database) read(w http.ResponseWriter, req *http.Request) {
	db.list(w, req)
}

func (db *database) update(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	newPrice := req.URL.Query().Get("price")

	price, err := strconv.ParseFloat(newPrice, 32)
	// Handle Parsing Failure
	if err != nil {
		w.WriteHeader(http.StatusBadRequest) // 400
		fmt.Fprintf(w, "invalid price: %q\n", newPrice)
		return
	}

	// Lock for Reading
	db.mu.RLock()
	_, ok := db.data[item]
	db.mu.RUnlock()

	if !ok {
		w.WriteHeader(http.StatusBadRequest) // 400
		fmt.Fprintf(w, "item does not exist: %s\n", item)
		return
	}

	db.mu.Lock()
	db.data[item] = dollars(price)
	db.mu.Unlock()

	fmt.Fprintf(w, "update item: %s, price: %s\n", item, db.data[item])
}

func (db *database) delete(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")

	// Lock and Key for Synchronization
	db.mu.RLock()
	_, ok := db.data[item]
	db.mu.RUnlock()

	if !ok {
		w.WriteHeader(http.StatusBadRequest) // 400
		fmt.Fprintf(w, "item does not exist: %s\n", item)
		return
	}

	db.mu.Lock()
	delete(db.data, item)
	db.mu.Unlock()

	fmt.Fprintf(w, "delete item: %s\n", item)
}
