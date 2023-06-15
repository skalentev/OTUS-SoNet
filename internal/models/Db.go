package models

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	//	_ "github.com/jackc/pgx"
	_ "github.com/lib/pq"
	"time"
	//	"gorm.io/driver/postgres"
	//	"gorm.io/gorm"
)

// var DB *gorm.DB

type Db struct {
	DB        *sql.DB
	Driver    string
	lastQuery string
}

var DB = &Db{}
var DBSlave = &Db{}

func (d *Db) Init(cfg DBConfig) error {
	var (
		err error
		dsn string
	)
	//	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port, cfg.SSLMode)
	//	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	d.Driver = cfg.Driver
	switch d.Driver {
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?interpolateParams=true", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	case "postgres":
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port, cfg.SSLMode)
	default:
		return errors.New("cannot open DB")
	}
	d.DB, err = sql.Open(d.Driver, dsn)
	if err != nil {
		return err
	}
	d.DB.SetConnMaxLifetime(time.Minute * 1)
	d.DB.SetMaxOpenConns(50)
	d.DB.SetMaxIdleConns(10)

	pingErr := d.DB.Ping()
	if pingErr != nil {
		return pingErr
	}
	//	fmt.Println("DB Connected!")

	switch d.Driver {
	case "mysql":
		d.lastQuery = `
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
		d.lastQuery = `
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

	_, err = d.DB.Exec(d.lastQuery)
	if err != nil {
		fmt.Println("Table user:", err)
	}

	switch d.Driver {
	case "mysql":
		d.lastQuery = `
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
		d.lastQuery = `
			CREATE TABLE IF NOT EXISTS public.session
			(
				token VARCHAR(64) NOT NULL,
				user_id VARCHAR(64) NOT NULL,
				created_at timestamp DEFAULT CURRENT_TIMESTAMP,
				updated_at timestamp DEFAULT CURRENT_TIMESTAMP,
				deleted_at timestamp,
				token_till timestamp
			);
			`
	}

	_, err = d.DB.Exec(d.lastQuery)
	if err != nil {
		fmt.Println("Table session:", err)
	}

	switch d.Driver {
	case "mysql":
		d.lastQuery = `
			CREATE TABLE IF NOT EXISTS post (
				id VARCHAR(64) NOT NULL PRIMARY KEY,
				user_id VARCHAR(64) NOT NULL,
				text TEXT,
				created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
				deleted_at TIMESTAMP,
				FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE
			) ENGINE=InnoDB
		`
	default:
		d.lastQuery = `
			CREATE TABLE IF NOT EXISTS public.post
			(
				id VARCHAR(64) NOT NULL,
				user_id VARCHAR(64) NOT NULL,
				text TEXT,
				created_at timestamp DEFAULT CURRENT_TIMESTAMP,
				updated_at timestamp DEFAULT CURRENT_TIMESTAMP,
				deleted_at timestamp
			);
			`
	}

	_, err = d.DB.Exec(d.lastQuery)
	if err != nil {
		fmt.Println("Table post:", err)
	}

	switch d.Driver {
	case "mysql":
		d.lastQuery = `
			CREATE TABLE IF NOT EXISTS friend (
				id VARCHAR(64) NOT NULL PRIMARY KEY,
				user_id VARCHAR(64) NOT NULL,
				friend_id VARCHAR(64) NOT NULL,
				created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
				deleted_at TIMESTAMP
			) ENGINE=InnoDB
		`
	default:
		d.lastQuery = `
			CREATE TABLE IF NOT EXISTS public.friend
			(
				id VARCHAR(64) NOT NULL,
				user_id VARCHAR(64) NOT NULL,
				friend_id VARCHAR(64) NOT NULL,
				created_at timestamp DEFAULT CURRENT_TIMESTAMP,
				updated_at timestamp DEFAULT CURRENT_TIMESTAMP,
				deleted_at timestamp
			);
			`
	}

	_, err = d.DB.Exec(d.lastQuery)
	if err != nil {
		fmt.Println("Table friend:", err)
	}

	return nil
}

func (d *Db) Close() error {
	return d.DB.Close()
}
