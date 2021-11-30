package db

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
	log "github.com/sirupsen/logrus"
)

func Connect() *sql.DB {
	var err error

	const steps = 2

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

	driver, _ := mysql.WithInstance(database, &mysql.Config{})
	mgt, _ := migrate.NewWithDatabaseInstance(
		"file://ConnectionDatabase/migrations",
		"mysql",
		driver,
	)

	mgt.Steps(steps)

	database.SetConnMaxLifetime(lifeTime * time.Minute)
	database.SetMaxOpenConns(openConns)
	database.SetMaxIdleConns(idleConns)

	err = database.Ping()
	if err != nil {
		log.Error("DB.Ping = ", err)
	}

	return database
}
