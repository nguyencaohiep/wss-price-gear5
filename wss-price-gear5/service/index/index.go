package index

import (
	"crawl_price_3rd/pkg/db"
	"crawl_price_3rd/pkg/router"
	"crawl_price_3rd/pkg/server"
	"net/http"
	"strings"
)

// GetIndex Function to Show API Information
func GetIndex(w http.ResponseWriter, r *http.Request) {
	router.ResponseSuccess(w, "200", "Crawler TOP COIN REALTIME is running")
}

// GetHealth Function to Show Health Check Status
func GetHealth(w http.ResponseWriter, r *http.Request) {
	// Check Database Connections
	if len(server.Config.GetString("DB_DRIVER")) != 0 {
		switch strings.ToLower(server.Config.GetString("DB_DRIVER")) {
		case "postgres":
			err := db.PSQL.Ping()
			if err != nil {
				router.ResponseInternalError(w, "postgres-health-check", err)
				return
			}
		}
	}

	// Return Success response
	router.ResponseSuccess(w, "200", "Health is ok")
}
