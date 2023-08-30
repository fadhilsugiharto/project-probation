package products

import (
	"database/sql"
	"log"
	"project-probation/model"
)

const (
	getQuery = "SELECT id, name, price FROM products WHERE CAST(id as TEXT) ILIKE '%' || $1 || '%' AND name ILIKE '%' || $2 || '%' AND CAST(price as TEXT) ILIKE '%' || $3 || '%'"

	insertQuery = "INSERT INTO products (name, price) VALUES ($1, $2) RETURNING id"
)

func Get(db *sql.DB, id string, name string, price string) ([]model.Product, error) {
	rows, err := db.Query(getQuery, id, name, price)
	if err != nil {
		log.Fatal("Error query to DB")
		return nil, err
	}

	defer rows.Close()

	products := []model.Product{}
	for rows.Next() {
		product := model.Product{}
		err := rows.Scan(&product.ID, &product.Name, &product.Price)
		if err != nil {
			log.Fatal("Error parsing data")
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
		log.Fatal("Error insert data")
		return insertId, err
	}

	return insertId, nil
}
