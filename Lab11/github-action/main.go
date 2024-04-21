// Launch microservice server- main.go
package main

import (
	"github.com/danielmyers1993/CloudNativeCourse"
	"log"
)

func main() {
	s := microservice.NewServer("", "8000")
	log.Fatal(s.ListenAndServe())
}
