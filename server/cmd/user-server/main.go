package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/kelseyhightower/envconfig"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	manager "github.com/leandro-lugaresi/grpc-realtime-chat/server/mysql"
	"github.com/leandro-lugaresi/grpc-realtime-chat/server/user"
	pb "github.com/leandro-lugaresi/grpc-realtime-chat/server/user/userpb"
)

var db *sql.DB

type config struct {
	DBHost          string
	DBPort          string `default:"3306"`
	DBUsername      string `default:"user-service"`
	DBPassword      string
	DBName          string `default:"chat"`
	TLSCert         string `default:"/etc/auth/cert.pem"`
	TLSKey          string `default:"/etc/auth/key.pem"`
	JWTPrivateKey   string `default:"/etc/auth/jwt-key.pem"`
	MaxAttempts     int    `default:"20"`
	ListenAddr      string `default:"127.0.0.1:7400"`
	DebugListenAddr string `default:"127.0.0.1:7401"`
}

func main() {
	var c config
	err := envconfig.Process("USER", &c)
	if err != nil {
		log.Fatalf("[CRITICAL][user-server] Could not process the config enviromment: %v \n", err)
	}
	log.Println("[INFO][user-server] User service starting...")

	var dbError error
	// Connect to database.
	db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", c.DBUsername, c.DBPassword, c.DBHost, c.DBPort, c.DBName))
	if err != nil {
		log.Printf("[ERROR][user-server] %v", err)
	}
	for attempts := 1; attempts < c.MaxAttempts; attempts++ {
		dbError = db.Ping()
		if dbError == nil {
			break
		}
		log.Printf("[ERROR][user-server] %v", dbError)
		time.Sleep(time.Duration(attempts) * time.Second)
	}
	if dbError != nil {
		log.Fatalf("[CRITICAL][user-server] Could not Connect to dababase: %v \n", dbError)
	}
	uM := manager.UserManager{db}

	ta, err := credentials.NewServerTLSFromFile(c.TLSCert, c.TLSKey)
	if err != nil {
		log.Fatalf("[CRITICAL][user-server] %v", err)
	}
	gs := grpc.NewServer(grpc.Creds(ta))

	key, err := ioutil.ReadFile(c.JWTPrivateKey)
	if err != nil {
		log.Fatalf("[CRITICAL][user-server] Error reading the jwt private key: %s", err)
	}
	us, err := user.NewUserService(key, uM)
	pb.RegisterUserServiceServer(gs, us)

	ln, err := net.Listen("tcp", c.ListenAddr)
	if err != nil {
		log.Fatalf("[CRITICAL][user-server] %v", err)
	}
	go gs.Serve(ln)

	log.Println("[INFO][user-server] User service started successfully.")
	log.Fatal(http.ListenAndServe(c.DebugListenAddr, nil))
}
