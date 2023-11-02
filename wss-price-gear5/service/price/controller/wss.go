package controller

import (
	"crawl_price_3rd/pkg/log"
	"crawl_price_3rd/service/price/crawler"
	"crawl_price_3rd/service/price/dao"
	"errors"
	"net/http"
	"sort"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
)

// ClientList is a map used to manage a map of clients
type ClientList map[*Client]bool

type Client struct {
	Connection *websocket.Conn
	// egress is used to avoid concurrent writes on the WebSocket
	ListPrice chan []dao.Crypto
}

func NewClient(conn *websocket.Conn) *Client {
	return &Client{
		Connection: conn,
		ListPrice:  make(chan []dao.Crypto),
	}
}

var (
	Clients ClientList
	Event   chan []dao.Crypto
)

func init() {
	Clients = ClientList{}
	Event = make(chan []dao.Crypto)
	go ListEvent()
	go BoardCast()
}

func HandleWS(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true } // check the http request

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(log.LogLevelError, "HandleWS upgrader.Upgrade(w, r, nil)", err.Error())
		return
	}

	// Create New Client
	client := NewClient(conn)
	// Add the newly created client to the manager
	Clients[client] = true
}

func BoardCast() {
	for {
		// time.Sleep(2 * time.Second)
		// fmt.Println("len clients", len(Clients))
		listPrice, hasPrice := <-Event
		if hasPrice {
			for clientEle := range Clients {
				err := clientEle.Connection.WriteJSON(listPrice)
				if err != nil {
					if !errors.Is(err, syscall.EPIPE) { // check err not broken pipe; broken pipe happen when client disconnect	ws
						log.Println(log.LogLevelError, "BoardCast clientEle.WriteJSON(listPrice)", err.Error())
					}
					clientEle.Connection.Close()
					delete(Clients, clientEle)
				}
			}
		} else {
			message := "ping"
			for clientEle := range Clients {
				err := clientEle.Connection.WriteJSON(message)
				if err != nil {
					if !errors.Is(err, syscall.EPIPE) {
						log.Println(log.LogLevelError, "BoardCast clientEle.WriteJSON(listPrice)", err.Error())
					}
					clientEle.Connection.Close()
					delete(Clients, clientEle)
				}
			}
		}
	}
}

func ListEvent() {
	for {
		time.Sleep(1 * time.Second)
		arrayTopPrice := crawler.ArrayBinance.Cryptos
		sort.Slice(arrayTopPrice, func(i, j int) bool {
			return arrayTopPrice[i].MarketcapUSD > arrayTopPrice[j].MarketcapUSD
		})

		resArr := []dao.Crypto{}
		for i, crypto := range arrayTopPrice {
			if i < 10 {
				resArr = append(resArr, crypto)
			}
		}
		Event <- resArr
	}
}
