package products

import (
	"database/sql"
	"log"
	"project-probation/model"
)

const (
	getQuery = "SELECT id, name, price FROM products"

	insertQuery = "INSERT INTO products (name, price) VALUES ($1, $2) RETURNING id"
)

func Get(db *sql.DB) ([]model.Product, error) {
	rows, err := db.Query(getQuery)
	if err != nil {
		log.Println("Error query to DB")
		return nil, err
	}

	defer rows.Close()

	products := []model.Product{}
	for rows.Next() {
		product := model.Product{}
		err := rows.Scan(&product.ID, &product.Name, &product.Price)
		if err != nil {
			log.Println("Error parsing data")
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

func Insert(db *sql.DB, name string, price string) (int, error) {
	var insertId int
	err := db.QueryRow(insertQuery, name, price).Scan(&insertId)
	if err != nil {
		log.Println("Error insert data")
		return insertId, err
	}

	return insertId, nil
}
