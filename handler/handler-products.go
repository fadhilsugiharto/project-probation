package handler

import (
	"context"
	"database/sql"
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

var isOutdated = true

type Handler interface {
	InsertProduct(w http.ResponseWriter, r *http.Request)
	GetProducts(w http.ResponseWriter, r *http.Request)
}

type handlerObj struct {
	dbCon *sql.DB
}

func New(db *sql.DB) Handler {
	return &handlerObj{
		dbCon: db,
	}
}

func (h *handlerObj) InsertProduct(w http.ResponseWriter, r *http.Request) {
	// Initiate db connection
	db, err := database.ConnectDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer db.Close()

	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	price := r.FormValue("price")

	// Insert product to db
	insertId, err := products.Insert(db, name, price)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Render the data in redis outdated
	isOutdated = true

	// Publish message to NSQ
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
	fmt.Fprintf(w, "Product successfully inserted!")
}

func (h *handlerObj) GetProducts(w http.ResponseWriter, r *http.Request) {
	var ctx = context.Background()

	//Try to get data from redis first
	if !isOutdated {
		productsFromRedis, err := redis.GetProductsFromRedis(ctx, "products")
		if err == nil {
			log.Println("Product data successfully retrieved from cache!")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(productsFromRedis)
			return
		}
	}

	// If data not exist, get from db
	// Initiate db connection
	db, err := database.ConnectDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer db.Close()

	// Get data from db
	res, err := products.Get(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Store data in Redis for future use
	err = redis.StoreProductsInRedis(ctx, "products", res)
	if err != nil {
		log.Println("Insert redis error")
	} else {
		isOutdated = false
		log.Println("Product successfully retrieved and cached!")
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
