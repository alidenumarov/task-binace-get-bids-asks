package handler

import (
	"binance-get-order-book/app"
	"binance-get-order-book/client/binance"
	"binance-get-order-book/models"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func EndpointGetBidsAndAsks(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	// определяем вебсокет соединение
	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Client Connected") // логируем соединение клиента через вебсокет

	basicUrl := "https://api.binance.com"
	httpClient := app.CreateBinanceHttpClient() // определяем клиента

	binanceClient := binance.NewBinanceClient(basicUrl, httpClient)

	for {
		bidsAndAsks, err := binanceClient.GetBidsAndAsks() // получение bids and asks
		if err != nil {
			wsRespMsg(wsConn, err.Error())
		}

		// высчитываем сумму объемов ордеров для BID
		bidsSum, err := calculateOrderBookVolume(bidsAndAsks.Bids)
		if err != nil {
			wsRespMsg(wsConn, err.Error())
		}

		// высчитываем сумму объемов ордеров для ASK
		asksSum, err := calculateOrderBookVolume(bidsAndAsks.Bids)
		if err != nil {
			wsRespMsg(wsConn, err.Error())
		}

		orderVolume := models.OrderVolume{
			Bids: *bidsSum,
			Asks: *asksSum,
		}

		resMsgByte, err := json.Marshal(orderVolume)
		if err != nil {
			wsRespMsg(wsConn, err.Error())
		}

		// отправляем результат на Веб
		if err := wsConn.WriteMessage(1, []byte(string(resMsgByte))); err != nil {
			log.Println(err)
			return
		}

		time.Sleep(1 * time.Second)
	}
}

func calculateOrderBookVolume(orderBook [][]string) (*float64, error) {
	var sum float64 // для объема ордеров в стакане BID или ASK
	items := 15     // Кол-во ордеров в каждом стакане - 15

	for i := 0; i < items; i++ {
		if len(orderBook[i]) < 2 {
			return nil, errors.New("wrong format returned from binance api, please, check documentation")
		}
		price, err := strconv.ParseFloat(orderBook[i][0], 64)
		if err != nil {
			return nil, err
		}
		qnt, err := strconv.ParseFloat(orderBook[i][1], 64)
		if err != nil {
			return nil, err
		}

		sum += price * qnt
	}

	return &sum, nil
}

// отправляем сообщение на Веб (независимо от результат это или ошибка)
func wsRespMsg(wsConn *websocket.Conn, message string) {
	if err := wsConn.WriteMessage(1, []byte(message)); err != nil {
		log.Println(err)
		return
	}
}
