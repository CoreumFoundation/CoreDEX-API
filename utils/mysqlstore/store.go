package store

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"

	logger "github.com/CoreumFoundation/CoreDEX-API/utils/logger"
)

type StoreBase struct {
	Client *sql.DB
}

type MySQLConfig struct {
	Username string
	Password string
	Host     string
	Port     int
	Database string
}

func Client() *StoreBase {
	// Parse the env variable MYSQL_CONFIG json (based on the MySQLConfig struct)
	e := os.Getenv("MYSQL_CONFIG")
	if e == "" {
		logger.Fatalf("MYSQL_CONFIG env variable is not set")
	}
	d := &MySQLConfig{}
	err := json.Unmarshal([]byte(e), d)
	if err != nil {
		logger.Fatalf("Error parsing MYSQL_CONFIG (%s) env variable: %v", e, err)
	}

	// Define the MySQL connection parameters
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		d.Username, d.Password, d.Host, d.Port, d.Database)

	// Open a connection to the MySQL server
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		logger.Errorf("Error opening database: %v", err)
		return nil
	}

	// Ping the database to verify the connection
	err = db.Ping()
	if err != nil {
		logger.Errorf("Error pinging database: %v", err)
		return nil
	}
	logger.Infof("Successfully connected to the MySQL database!")
	return &StoreBase{Client: db}
}
