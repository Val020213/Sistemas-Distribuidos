package main

import (
	"fmt"
	"log"
	"net/http"
)

const portNum string = ":8080"

func getHome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Homepage")
}

func getInfo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Info page")
}

func main() {
	log.Println("Starting http server")

	http.HandleFunc("/", getHome)
	http.HandleFunc("/info", getInfo)

	log.Println("Started on port", portNum)

	err := http.ListenAndServe(portNum, nil)
	if err != nil {
		log.Fatal(err)
	}
}
