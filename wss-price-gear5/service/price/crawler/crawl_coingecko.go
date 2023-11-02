package crawler

import (
	"crawl_price_3rd/pkg/log"
	"crawl_price_3rd/pkg/utils"
	"crawl_price_3rd/service/price/dao"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type CoinGeckoMarketInfo struct {
	ID                    string  `json:"id"`
	CurrentPrice          float64 `json:"current_price"`
	MarketCap             float64 `json:"market_cap"`
	TotalSupply           float64 `json:"total_supply"`
	PricePercentChange24h float64 `json:"price_change_percentage_24h"`
}

var clientCoingecko http.Client

func init() {
	clientCoingecko = http.Client{}
}

func CrawlPriceCoingecko() {
	page := 1

	for lenListPerPage := -1; lenListPerPage != 0 && page <= 50; page++ {
		repo := dao.ListCrypto{
			CryptoSrc: SrcCGC,
		}
		api := fmt.Sprintf(`https://api.coingecko.com/api/v3/coins/markets?vs_currency=usd&per_page=%v&page=%v`, 250, page)

		resp, err := clientCoingecko.Get(api)
		if err != nil {
			log.Println(log.LogLevelWarn, "CrawlPriceCoingecko clientCoingecko.Get(coingeckoAPI)", err.Error())
			return
		}

		if resp != nil {
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Println(log.LogLevelWarn, "Coingecko/CrawlPrices", err.Error())
				return
			}
			defer resp.Body.Close()

			rawCoingeckoDTOArr := make([]any, 0)
			err = json.Unmarshal(body, &rawCoingeckoDTOArr)
			if err != nil {
				log.Println(log.LogLevelWarn, "CrawlPriceCoingecko Unmarshal(body, &rawCoingeckoDTOArr)"+string(body), err.Error())
				return
			}
			lenListPerPage = len(rawCoingeckoDTOArr)

			// Traverse each json object from response array data got above.
			for _, rawCoingeckoDTO := range rawCoingeckoDTOArr {
				coinGeckoMarketInfo := &CoinGeckoMarketInfo{}
				err = utils.Mapping(rawCoingeckoDTO, coinGeckoMarketInfo)
				if err != nil {
					log.Println(log.LogLevelWarn, "CrawlPriceCoingecko utils.Mapping(rawCoingeckoDTO, coingeckoDTO)", err.Error())
					continue
				}

				mutex.Lock()
				crypto, exist := MapCryptocodeCGC[coinGeckoMarketInfo.ID]

				if exist {
					crypto := &dao.Crypto{
						CryptoId:              crypto.CryptoId,
						PriceUSD:              coinGeckoMarketInfo.CurrentPrice,
						Symbol:                crypto.Symbol,
						MarketcapUSD:          coinGeckoMarketInfo.MarketCap,
						PricePercentChange24h: coinGeckoMarketInfo.PricePercentChange24h,
						Name:                  crypto.Name,
					}
					MapCryptocodeCGC[coinGeckoMarketInfo.ID] = *crypto
					repo.Cryptos = append(repo.Cryptos, *crypto)
				}

				mutex.Unlock()
			}
		}
		if page%2 == 0 {
			time.Sleep(30 * time.Second)
		}
	}
}
