package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv" //imported to parse string value of price into a float32 value
	"sync"    //provides mutual exlusion in order to sync across database
)

func main() {
	db := database{"shoes": 50, "socks": 5}
	http.HandleFunc("/list", db.list) //handlers for list, price, create, update, delete so it can read address properly
	http.HandleFunc("/price", db.price)
	http.HandleFunc("/create", db.create)
	http.HandleFunc("/update", db.update)
	http.HandleFunc("/delete", db.delete)
	log.Fatal(http.ListenAndServe("localhost:8000", nil)) //no argmuent needs to be passed so nil is passed
}

type dollars float32 // Dollars will be float32 to accommodate 32 bit decimals

func (d dollars) String() string { return fmt.Sprintf("$%.2f", d) }

type database map[string]dollars

func (db database) list(w http.ResponseWriter, req *http.Request) {
	for item, price := range db {
		fmt.Fprintf(w, "%s: %s\n", item, price)
	}
}

func (db database) price(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	price, ok := db[item]
	if !ok {
		w.WriteHeader(http.StatusNotFound) // 404
		fmt.Fprintf(w, "no such item: %q\n", item)
		return
	}
	fmt.Fprintf(w, "%s\n", price)
}

// Create
func (db database) create(w http.ResponseWriter, req *http.Request) { //takes in 2 arguments W, type http.responsewriter, req type *httprequest
	//w used as instance to construct http response, passed as parameter to http handler func so that func can write reponse back to client
	//req passed as parameter to function so that it can access information about type of request made, use to give appropriate response
	query := req.URL.Query()    //grab query parameters of url request
	item := query.Get("item")   //capture item name
	price := query.Get("price") //get price value

	// Check if item and price are present in the request
	if item == "" || price == "" { //if item and price are empty
		w.WriteHeader(http.StatusNotFound)          //updates status of error by writing price not found in w ResponseWriter
		fmt.Fprintf(w, "item or price not found\n") //prints error msg
		return
	}

	// Convert the price string to float32
	p, err := strconv.ParseFloat(price, 32) //convert price string to float32 value
	if err != nil {                         //if there is an error/fails to convert (as indicated by err does not equal to nil)
		w.WriteHeader(http.StatusNotFound) //updates w ResponseWriter with invalid price status
		fmt.Fprintf(w, "invalid price\n")
		return
	}

	// Use a Mutex to synchronize access to the db
	var mutex = &sync.Mutex{} //new mutex object to sync database
	mutex.Lock()              //lock database before updating
	db[item] = dollars(p)     //update database(adds new item with price)
	mutex.Unlock()            //unlock after update

	fmt.Fprintf(w, "item %s with price %s successfully created\n", item, price) //print successful message
}

// Update
func (db database) update(w http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()    //grab query parameters of url request
	item := query.Get("item")   //get item name
	price := query.Get("price") //get item price

	// Check if item and price are present in the request
	if item == "" || price == "" {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "item or price not found\n")
		return
	}

	// Convert the price string to float32
	p, err := strconv.ParseFloat(price, 32)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "invalid price\n")
		return
	}

	// Check if the item exists in the db
	_, ok := db[item] //
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "no such item: %q\n", item)
		return
	}

	// Use a Mutex to synchronize access to the db
	var mutex = &sync.Mutex{}
	mutex.Lock()
	db[item] = dollars(p)
	mutex.Unlock()

	fmt.Fprintf(w, "item %s with price %s successfully updated\n", item, price)
}

// Delete
func (db database) delete(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item") //get url

	// Check if item is present in the request
	if item == "" { //if item field is empty, update w responsewriter with error
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "item not found\n")
		return
	}

	// Check if the item exists in the db
	_, ok := db[item] //use item as key to look for in map, ok indicates whether or not item is in there
	if !ok {          //if it is not in there/key is false
		w.WriteHeader(http.StatusNotFound)         //display error
		fmt.Fprintf(w, "no such item: %q\n", item) // writes err msg to response writer
		return
	}

	// Use a Mutex to synchronize access to the db
	var mutex = &sync.Mutex{} //new instance of mutex for sync
	mutex.Lock()              //locks db
	delete(db, item)          //update db
	mutex.Unlock()            //unlock db again

	fmt.Fprintf(w, "item %s successfully deleted\n", item) //print successfull message
}
