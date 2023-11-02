package crawler

import (
	"bytes"
	"crawl_price_3rd/pkg/log"
	"crawl_price_3rd/pkg/server"
	"crawl_price_3rd/service/price/dao"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"sync"
)

var (
	MapCryptocodeCGC map[string]dao.Crypto
	MapCryptocodeCMC map[string]string
	MapCryptocodeSol map[string]string
	MapPriceBinance  map[string]string
	MapCryptoCode    map[string]string // to use link from MapPriceBinance to MapCryptocodeCGC to get crypto info
	NumberUpdateCMC  int
	NumberUpdateSol  int
	mutex            sync.Mutex
	LenCryptosCMC    int
	LenCryptosSol    int
)

var (
	SrcCGC = "coingecko"
	SrcCMC = "coinmarketcap"
	SrcSOL = "solscan"
	SrcBNB = "binance"
)

type InfoUpdate struct {
	LastUpdateTime string `json:"LastUpdateTime"`
	Update         int    `json:"Update"`
	Insert         int    `json:"Insert"`
}

type DataInfo struct {
	Status  bool   `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    Data   `json:"data"`
}
type Data struct {
	Infos []Info `json:"infos"`
}
type Info struct {
	Cryptoid   string `json:"cryptoid"`
	Cryptosrc  string `json:"cryptosrc"`
	Cryptocode string `json:"cryptocode"`
	Symbol     string `json:"symbol"`
	Name       string `json:"name"`
}

func init() {
	MapCryptocodeCGC = map[string]dao.Crypto{}
	MapPriceBinance = map[string]string{}
	prepareInfo()
}

func prepareInfo() error {
	api := server.Config.GetString("DOMAIN_LOCAL") + server.Config.GetString("API_GET_INFO") // call api get info from db
	req, err := http.NewRequest(http.MethodGet, api, nil)
	if err != nil {
		log.Println(log.LogLevelError, "prepareInfo", err.Error())
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Println(log.LogLevelError, "client.Do(req)", err.Error())
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(log.LogLevelError, "io.ReadAll(res.Body)", err.Error())
		return err
	}

	infoData := &DataInfo{}
	err = json.Unmarshal(body, &infoData)
	if err != nil {
		log.Println(log.LogLevelError, " json.Unmarshal(body, &infoData)", err.Error())
		return err
	}

	numberTop := 0

	for _, info := range infoData.Data.Infos {
		crypto := dao.Crypto{
			Cryptocode: info.Cryptocode,
			CryptoId:   info.Cryptoid,
			CryptoSrc:  info.Cryptosrc,
			Symbol:     info.Symbol,
			Name:       info.Name,
		}
		if info.Cryptosrc == SrcCGC {
			_, exist := MapCryptocodeCGC[info.Cryptocode]
			if !exist {
				MapCryptocodeCGC[info.Cryptocode] = crypto
				if numberTop < 100 {
					numberTop++
					_, exist := MapPriceBinance[crypto.Symbol+"USDT"]
					if !exist {
						MapPriceBinance[crypto.Symbol+"USDT"] = crypto.Cryptocode
					}
				}
			}
		}
	}
	return nil
}

func UpdatePrice(repo dao.ListCrypto) error {
	// return nil
	jsonBody, err := json.Marshal(repo)
	if err != nil {
		return err
	}
	bodyReader := bytes.NewReader(jsonBody)

	api := server.Config.GetString("DOMAIN_LOCAL") + server.Config.GetString("API_UPDATE_PRICE")
	req, err := http.NewRequest(http.MethodPatch, api, bodyReader)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode > 299 || res.StatusCode < 200 {
		return errors.New(res.Status + " " + string(resBody))
	}
	return nil
}
