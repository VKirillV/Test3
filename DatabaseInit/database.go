package init

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)


func init() {
	
	e := godotenv.Load()
	if e != nil {
		log.Println(e)
	}

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

	//overtime times
	DB.SetConnMaxLifetime(3 * time.Minute)
	// maximum connection number
	DB.SetMaxOpenConns(100)
	// Set the number of idle connections
	DB.SetMaxIdleConns(16)
	if err := DB.Ping(); err != nil {
		log.Error("DB.Ping = ", err)
	}

}
