package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"project-probation/database"
	"project-probation/database/products"
	"project-probation/model"
	"project-probation/redis"
	"strconv"
)

func InsertProduct(w http.ResponseWriter, r *http.Request) {
	db, err := database.ConnectDB()
	if err != nil {
		http.Error(w, "Failed to connect to the database", http.StatusInternalServerError)
		return
	}

	defer db.Close()

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	price := r.FormValue("price")

	insertId, err := products.Insert(db, name, price)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	priceInt, err := strconv.Atoi(price)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	product := model.Product{
		ID:    insertId,
		Name:  name,
		Price: priceInt,
	}

	productJSON, _ := json.Marshal(product)

	err = nsqProducer.Publish("project_probation", productJSON)
	if err != nil {
		log.Println("Failed to publish message to NSQ")
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Product inserted successfully!")
}

func GetProducts(w http.ResponseWriter, r *http.Request) {
	var ctx = context.Background()

	id := r.FormValue("id")
	name := r.FormValue("name")
	price := r.FormValue("price")
	redisKey := fmt.Sprintf("id:%s,name:%s,price:%s", id, name, price)

	//Try to get data from redis first
	productsFromRedis, err := redis.GetProductsFromRedis(ctx, redisKey)
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(productsFromRedis)
		return
	}

	// If data not exist, get from db
	// Initiate db connection
	db, err := database.ConnectDB()
	if err != nil {
		http.Error(w, "Failed to connect to the database", http.StatusInternalServerError)
		return
	}

	defer db.Close()

	// Get data from db
	res, err := products.Get(db, id, name, price)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Store data in Redis for future use
	err = redis.StoreProductsInRedis(ctx, redisKey, res)
	if err != nil {
		log.Println("Insert redis error")
	} else {
		log.Println("products data is cached")
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
