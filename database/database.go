package database

import (
	"database/sql"
	"fmt"
	"project-probation/consts"

	_ "github.com/lib/pq"
)

func ConnectDB() (*sql.DB, error) {

	postgres_url := fmt.Sprintf(
		"%s://%s:%s@%s:%s/%s?sslmode=disable",

		consts.Db,
		consts.User,
		consts.Password,
		consts.Host,
		consts.Port,
		consts.Dbname,
	)

	return sql.Open("postgres", postgres_url)
}
