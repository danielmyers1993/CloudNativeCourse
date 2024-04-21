// Launch microservice server- main.go
package main

import (
	"https://github.com/danielmyers1993/CloudNativeCourse/tree/main/Lab11/github-action/microservice"
	"log"
)

func main() {
	s := microservice.NewServer("", "8000")
	log.Fatal(s.ListenAndServe())
}
