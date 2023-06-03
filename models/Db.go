package models

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	//	_ "github.com/jackc/pgx"
	_ "github.com/lib/pq"
	"log"
	"time"
	//	"gorm.io/driver/postgres"
	//	"gorm.io/gorm"
)

// var DB *gorm.DB

var DB *sql.DB
var Driver string

func InitDB(cfg Config) {

	var (
		err error
		db  *sql.DB
		dsn string
	)
	//	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port, cfg.SSLMode)
	//	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	Driver = cfg.Driver
	switch Driver {
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?interpolateParams=true", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	case "postgres":
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port, cfg.SSLMode)
	default:
		panic("Cannot open DB: set .env Driver value as mysql or postgres")
	}
	db, err = sql.Open(Driver, dsn)

	if err != nil {
		panic(err)
	}
	db.SetConnMaxLifetime(time.Minute * 1)
	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(10)

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("DB Connected!")

	var initScript string
	switch Driver {
	case "mysql":
		initScript = `
			CREATE TABLE IF NOT EXISTS user (
				id VARCHAR(64) NOT NULL,
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
			)  ;
		`
	default:
		initScript = `
			CREATE TABLE IF NOT EXISTS public.user
			(
				id VARCHAR(64) NOT NULL,
				created_at timestamp DEFAULT CURRENT_TIMESTAMP,
				updated_at timestamp DEFAULT CURRENT_TIMESTAMP,
				deleted_at timestamp,
				first_name VARCHAR(90) NOT NULL,
				second_name VARCHAR(90) NOT NULL,
				birthdate VARCHAR(20) NOT NULL,
				biography text,
				city VARCHAR(64),
				password VARCHAR(64)
			);
		`
	}

	_, err = db.Exec(initScript)
	if err != nil {
		panic(err)
	}

	switch Driver {
	case "mysql":
		initScript = `
			CREATE TABLE IF NOT EXISTS session (
				token VARCHAR(64) NOT NULL PRIMARY KEY,
				user_id VARCHAR(64) NOT NULL,
				created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
				deleted_at TIMESTAMP,
				token_till TIMESTAMP,
				FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE
			) ENGINE=InnoDB
		`
	default:
		initScript = `
			CREATE TABLE IF NOT EXISTS public.session
			(
				token VARCHAR(64) NOT NULL,
				user_id VARCHAR(64) NOT NULL,
				created_at timestamp,
				updated_at timestamp,
				deleted_at timestamp,
				token_till timestamp
			);
			`
	}

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
