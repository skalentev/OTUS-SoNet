package models

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
	//	"gorm.io/driver/postgres"
	//	"gorm.io/gorm"
)

// var DB *gorm.DB

var DB *sql.DB

func InitDB(cfg Config) {

	var (
		err error
		db  *sql.DB
	)
	//	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port, cfg.SSLMode)
	//	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?interpolateParams=true", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(100)

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("DB Connected!")

	var initScript string = `
		CREATE TABLE IF NOT EXISTS user (
		    id VARCHAR(36) NOT NULL,
		    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		    deleted_at TIMESTAMP,
		    first_name VARCHAR(90) NOT NULL,
		    second_name VARCHAR(90) NOT NULL,
		    birthdate VARCHAR(20) NOT NULL,
		    biography TEXT,
		    city VARCHAR(64),
		    password VARCHAR(64),
		    PRIMARY KEY (id)
		) ENGINE=InnoDB ;
 `
	_, err = db.Exec(initScript)
	if err != nil {
		panic(err)
	}

	initScript = `CREATE TABLE IF NOT EXISTS session (
		token VARCHAR(64) NOT NULL PRIMARY KEY,
		user_id VARCHAR(64) NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		deleted_at TIMESTAMP,
		token_till TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE
	) ENGINE=InnoDB`
	_, err = db.Exec(initScript)
	if err != nil {
		panic(err)
	}
	fmt.Println("DB Initialized!")

	DB = db
}

func CloseDB() {
	err := DB.Close()
	if err != nil {
		panic(err)
	}
}
