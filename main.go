package main

import (
	"binance-get-order-book/router"
	"fmt"
	"log"
	"net/http"
)

func main() {
	router.SetupRoutes()

	fmt.Println("listening on port:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
