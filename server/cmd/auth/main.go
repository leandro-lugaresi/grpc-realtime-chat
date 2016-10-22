package main

import (
	"database/sql"
	"log"
	"net"
	"time"

	"github.com/kelseyhightower/envconfig"
)

var db *sql.DB

type Config struct {
	DBHost        string
	DBPort        string `default:"3306"`
	DBUsername    string `default:"auth"`
	DBPassword    string
	DBName        string `default:"chat"`
	TLSCert       string `default:"/etc/auth/cert.pem"`
	TLSKey        string `default:"/etc/auth/key.pem"`
	JWTPrivateKey string `default:"/etc/auth/jwt-key.pem"`
	MaxAttempts   int8   `default:"20"`
}

func main() {
	var c Config
	err := envconfig.Process("AUTH", &c)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println("Auth service starting...")

	var dbErr error
	// Connect to database.
	dbAddr := net.JoinHostPort(c.DBHost, c.DBPort)
	dbConfig := mysql.Config{
		User:   c.DBUsername,
		Passwd: c.DBPassword,
		Net:    "tcp",
		Addr:   dbAddr,
		DBName: c.DBName,
	}
	db, err = sql.Open("mysql", dbConfig.FormatDSN())
	if err != nil {
		log.Println(err)
	}
	for attempts := 1; attempts < c.MaxAttempts; attempts++ {
		dbError = db.Ping()
		if dbError == nil {
			break
		}
		log.Println(dbError)
		time.Sleep(attempts * time.Second)
	}
	if dbError != nil {
		log.Fatal(dbError)
	}
}
