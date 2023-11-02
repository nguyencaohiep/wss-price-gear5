package crawler

import (
	"crawl_price_3rd/pkg/log"
	"crawl_price_3rd/service/price/dao"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

var RepoInfoBinance []PriceBinance
var clientBinance http.Client
var ArrayBinance dao.ListCrypto

func init() {
	clientBinance = http.Client{}
	ArrayBinance = dao.ListCrypto{}
}

type PriceBinance struct {
	Symbol string `json:"symbol"` // symbol binanace return, ex BTCUSDT
	Price  string `json:"price"`
}

func CrawlPriceBinance() {
	repo := dao.ListCrypto{
		CryptoSrc: SrcBNB,
	}
	arrayBinance := dao.ListCrypto{}
	api := `https://api.binance.com/api/v3/ticker/price`

	resp, err := clientBinance.Get(api)
	if err != nil {
		log.Println(log.LogLevelWarn, "CrawlPriceBinance client.Get(api)", err.Error())
		return
	}

	if resp != nil {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(log.LogLevelWarn, "CrawlPriceBinance ioutil.ReadAll(resp.Body)", err.Error())
			return
		}
		defer resp.Body.Close()

		err = json.Unmarshal(body, &RepoInfoBinance)
		if err != nil {
			log.Println(log.LogLevelWarn, "CrawlPriceBinance json.Unmarshal(body, &resSol)", err.Error())
			return
		}

		for _, priceInfo := range RepoInfoBinance {
			cryptoCode, exist := MapPriceBinance[priceInfo.Symbol]
			if exist {
				priceFloat, err := strconv.ParseFloat(priceInfo.Price, 64)
				if err == nil {
					mutex.Lock()
					cryptoEle := MapCryptocodeCGC[cryptoCode]
					crypto := &dao.Crypto{ // to update db,
						CryptoId: cryptoEle.CryptoId,
						PriceUSD: priceFloat,
					}
					repo.Cryptos = append(repo.Cryptos, *crypto)

					cryptoEle.PriceUSD = priceFloat
					arrayBinance.Cryptos = append(arrayBinance.Cryptos, cryptoEle)
					mutex.Unlock()
				}
			}
		}
		ArrayBinance = arrayBinance
	}
	err = UpdatePrice(repo)
	if err != nil {
		log.Println(log.LogLevelError, "CrawlPriceBinance: UpdatePrice", err.Error())
	}
	// fmt.Println("len bnb", len(repo.Cryptos))
}
