package main

import (
	"fmt"
	"log"
	"net/http"
	"project-probation/handler"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	//go handler.ConsumeNSQMessages()
	r := mux.NewRouter()

	// handle get products
	r.HandleFunc("/products", handler.GetProducts).Methods("GET")

	// handle insert products
	r.HandleFunc("/products/insert", handler.InsertProduct).Methods("POST")

	log.Println("Server is running on :8080")
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
	fmt.Println(7)
}
