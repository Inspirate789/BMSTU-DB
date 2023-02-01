package database

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

// Execute from project root: go:generate go-bindata -pkg database -o database/sqlscripts.go ./database/sql

/*
	AuthDB Example: {
		host:    	"postgres"
		port:     	5432
		username: 	"user"
		password: 	"mypassword"
		dbname:   	"user_db"
		sslmode:  	"disable"
	}
*/
type AuthDB struct {
	Host     string
	Port     int
	Username string
	Password string
	DBName   string
	SSLMode  string
}

type Database struct {
	sqlDB       *sql.DB
	gormDB      *gorm.DB
	authData    AuthDB
	ConnTimeout time.Duration
}

func (d *Database) Connect(user AuthDB) error {
	postgresInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		user.Host, user.Port, user.Username, user.Password, user.DBName, user.SSLMode)
	sqlDB, err := sql.Open("postgres", postgresInfo)
	if err != nil {
		return err
	}
	d.sqlDB = sqlDB

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: d.sqlDB}), &gorm.Config{})
	if err != nil {
		return err
	}
	d.gormDB = gormDB

	d.authData = user

	end := time.Now().Add(d.ConnTimeout)
	for d.sqlDB.Ping() != nil {
		if time.Now().After(end) {
			return errors.New(fmt.Sprintf("failed to connect after %v secs", d.ConnTimeout.Seconds()))
		}
	}

	return nil
}

//func (d *Database) Init() error {
//	files, _ := filepath.Glob("database/sql/*_init.sql")
//
//	for _, file := range files {
//		script, err := Asset(file)
//		if err != nil {
//			return err
//		}
//
//		_, err = d.sqlDB.Query(string(script))
//		if err != nil {
//			return err
//		}
//	}
//
//	return nil
//}

func (d *Database) Disconnect() error {
	return d.sqlDB.Close()
}
