package dao

type ListCrypto struct {
	CryptoSrc string   `json:"cryptosrc"`
	Cryptos   []Crypto `json:"cryptos"`
}

type Crypto struct {
	CryptoId              string  `json:"cryptoid"`
	Name                  string  `json:"name"`
	CryptoSrc             string  `json:"cryptosrc"`
	Cryptocode            string  `json:"cryptocode"`
	Symbol                string  `json:"symbol"`
	Address               string  `json:"address"`
	MarketcapUSD          float64 `json:"marketcapusd"`
	TotalSupply           string  `json:"totalSupply"`
	PriceUSD              float64 `json:"priceUSD"`
	PricePercentChange24h float64 `json:"pricePercentChange24h"`
}
