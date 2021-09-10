package router

import (
	"binance-get-order-book/handler"
	"net/http"
)

func SetupRoutes() {
	http.HandleFunc("/health", handler.EndpointCheckHealth) // to check server's health
	http.HandleFunc("/ws", handler.EndpointGetBidsAndAsks)  // to get bids and asks volume
}
