package price

import (
	"crawl_price_3rd/service/price/controller"
	"crawl_price_3rd/service/price/crawler"
	"net/http"
	"time"

	"github.com/go-chi/chi"
)

var PriceServiceSunRouter = chi.NewRouter()

func init() {
	go func() {
		for {
			crawler.CrawlPriceCoingecko()
			// time.Sleep(5 * time.Minute)
		}
	}()

	go func() {
		for {
			crawler.CrawlPriceBinance()
			time.Sleep(time.Second)
		}
	}()

	PriceServiceSunRouter.Group(func(r chi.Router) {
		PriceServiceSunRouter.Handle("/crypto/top", http.HandlerFunc(controller.HandleWS))
	})
}
