package main

import (
	"log"
	"net/http"
	"project-probation/database"
	"project-probation/handler"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	go handler.ConsumeNSQMessages()
	r := mux.NewRouter()
	db, err := database.ConnectDB()
	if err != nil {
		log.Printf("err connect to db %v", err.Error())
		return
	}
	h := handler.New(db)

	// handle get products
	r.HandleFunc("/products", h.GetProducts).Methods("GET")
	// handle insert products
	r.HandleFunc("/products/insert", h.InsertProduct).Methods("POST")

	log.Println("Server is running on :8080")
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}
