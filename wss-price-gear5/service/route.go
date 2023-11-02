package service

import (
	"crawl_price_3rd/pkg/router"
	"crawl_price_3rd/service/index"
	"crawl_price_3rd/service/price"

	"github.com/go-chi/chi/middleware"
)

// LoadRoutes to Load Routes to Router
func LoadRoutes() {

	router.Router.Use(middleware.RealIP)

	//* Set Endpoint for admin
	router.Router.Get(router.RouterBasePath+"/", index.GetIndex)
	router.Router.Get(router.RouterBasePath+"/health", index.GetHealth)
	router.Router.Mount(router.RouterBasePath+"/prices", price.PriceServiceSunRouter)

}
