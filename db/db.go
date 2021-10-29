package db

import (
	"database/sql"
	"fmt"

	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
	log "github.com/sirupsen/logrus"
)

func InitDB() *sql.DB {

	var err error

	db_name := os.Getenv("db_name")
	db_host := os.Getenv("db_host")
	db_user := os.Getenv("db_user")
	db_pass := os.Getenv("db_pass")
	db_port := os.Getenv("db_port")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", string(db_user), string(db_pass), string(db_host), string(db_port), string(db_name))
	DB, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Panic("Failed to connect to database: ", err)
	}
	driver, _ := mysql.WithInstance(DB, &mysql.Config{})
	m, _ := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		"mysql",
		driver,
	)
	m.Steps(2)
	DB.SetConnMaxLifetime(3 * time.Minute)
	DB.SetMaxOpenConns(100)
	DB.SetMaxIdleConns(16)
	if err := DB.Ping(); err != nil {
		log.Error("DB.Ping = ", err)
	}
	return DB
}
