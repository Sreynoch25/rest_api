package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func InitDB() *sql.DB {

	dbInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	createTableSQL := `
CREATE TABLE IF NOT EXISTS tbl_users (
        id SERIAL PRIMARY KEY,
        last_name VARCHAR(30) NOT NULL,
        first_name VARCHAR(30) NOT NULL,
        user_name VARCHAR(50) NOT NULL UNIQUE,
        login_id VARCHAR(100) NOT NULL UNIQUE,
        email VARCHAR(50) NOT NULL,
        password VARCHAR(255) NOT NULL,
        role_name VARCHAR(50) NOT NULL,
        role_id INTEGER NOT NULL,
        is_admin BOOLEAN DEFAULT FALSE,
        login_session TEXT DEFAULT NULL,
        last_login TIMESTAMP WITHOUT TIME ZONE DEFAULT NULL,
        currency_id INTEGER,
        language_id INTEGER,
        status_id SMALLINT NOT NULL DEFAULT 1,
        "order" INTEGER,
        created_by INTEGER NOT NULL,
        created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_by INTEGER NOT NULL,
        updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
        deleted_by INTEGER,
        deleted_at TIMESTAMP WITHOUT TIME ZONE
    );
CREATE UNIQUE INDEX IF NOT EXISTS unique_email ON tbl_users (email) WHERE deleted_at IS NULL;
`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Database connected and initialized successfully")

	return db
}
