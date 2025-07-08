package services

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

func Services() {
	// Database connection details
	//connStr := "host=localhost port=5432 user=postgres password=replan dbname=replan sslmode=disable"
	connStr := "host=dpg-d1md11mmcj7s73a5qt50-a port=5432 user=replan_icjw_user password=klOhlKiWtnpI3vK14rkoMkgP9Sw2MUwR dbname=replan_icjw sslmode=require"
	// Open database connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error opening database:", err)
	}
	defer db.Close()

	// Check if connection is successful
	err = db.Ping()
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}

	fmt.Println("Successfully connected to the database!")
}
