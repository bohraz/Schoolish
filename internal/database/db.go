package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)


var (
	DB *sql.DB
	err error

	driverName string = "mysql"// database name
	dataSourceName string = "root:root@(127.0.0.1:3306)/mysql?parseTime=true" // database connection
)

// Initializes the database connection
func Init() {
	DB, err = sql.Open(driverName, dataSourceName)
	if err != nil {
		log.Fatalln("Error opening MySQL: ", err)
	}

	// Verifies the database connection
	err = DB.Ping()
	if err != nil {
		fmt.Println("Error connecting to the database: ", err)
	}
	
	fmt.Println("Successfully connected to the database!")
}

// Closes the database connection
func Close() {
	if DB != nil {
		err := DB.Close()
		if err != nil {
			fmt.Println("Successfully closed the database!")
		}
	}
}