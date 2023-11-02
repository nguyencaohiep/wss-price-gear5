package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"crawl_price_3rd/pkg/db"
	"crawl_price_3rd/pkg/router"
	"crawl_price_3rd/pkg/server"
	"crawl_price_3rd/service"
)

// Server Variable
var svr *server.Server

// Init Function
func init() {
	// Set Go Log Flags
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	// Load Routes
	service.LoadRoutes()

	// Initialize Server
	svr = server.NewServer(router.Router)
}

// Main Function
func main() {
	// Starting Server
	svr.Start()

	sig := make(chan os.Signal, 1)
	// Notify Any Signal to OS Signal Channel
	signal.Notify(sig, os.Interrupt)
	signal.Notify(sig, syscall.SIGTERM)

	// Return OS Signal Channel
	// As Exit Sign
	<-sig

	// Log Break Line
	fmt.Println("")

	// Stopping Server
	defer svr.Stop()

	// Close Any Database Connections
	if len(server.Config.GetString("DB_DRIVER")) != 0 {
		switch strings.ToLower(server.Config.GetString("DB_DRIVER")) {
		case "postgres":
			log.Println("Stoped postgres !")
			defer db.PSQL.Close()
		case "mysql":
			log.Println("Stoped mysql !")
			defer db.MySQL.Close()
		case "mongo":
			log.Println("Stoped mongo !")
			defer db.MongoSession.Close()
		}
	}
}
