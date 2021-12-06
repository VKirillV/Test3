package db

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

func Connect() *sql.DB {
	var err error

	const openConns = 100

	const idleConns = 16

	const lifeTime = 3

	dbName := os.Getenv("db_name")
	dbHost := os.Getenv("db_host")
	dbUser := os.Getenv("db_user")
	dbPass := os.Getenv("db_pass")
	dbPort := os.Getenv("db_port")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", dbUser, dbPass,
		dbHost, dbPort, dbName)

	database, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Panic("Failed to connect to database: ", err)
	}

	database.SetConnMaxLifetime(lifeTime * time.Minute)
	database.SetMaxOpenConns(openConns)
	database.SetMaxIdleConns(idleConns)

	err = database.Ping()
	if err != nil {
		log.Error("DB.Ping = ", err)
	}

	return database
}
