package db

import (
	"crawl_price_3rd/pkg/log"

	mgo "gopkg.in/mgo.v2"
)

// Mongo Configuration Struct
type mongoConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

// Mongo Configuration Variable
var mongoCfg mongoConfig

// MongoSession Variable
var MongoSession *mgo.Session

// Mongo Variable
var Mongo *mgo.Database

// Mongo Connect Function
func mongoConnect() (*mgo.Session, *mgo.Database) {
	// Initialize Connection
	conn, err := mgo.Dial(mongoCfg.User + ":" + mongoCfg.Password + "@" + mongoCfg.Host + ":" + mongoCfg.Port + "/" + mongoCfg.Name)
	if err != nil {
		log.Println(log.LogLevelFatal, "mongo-connect", err.Error())
	}

	// Test Connection
	err = conn.Ping()
	if err != nil {
		log.Println(log.LogLevelFatal, "mongo-connect", err.Error())
	} else {
		log.Println(log.LogLevelInfo, "mongo-connect", "Connect mongo: Successfully connected")
	}

	// Set Connection to Monotonic
	conn.SetMode(mgo.Monotonic, true)

	// Return Connection
	return conn, conn.DB(mongoCfg.Name)
}
