package app

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"
)

type AppContext struct {
}

func CreateBinanceHttpClient() *http.Client {
	cli := new(http.Client)
	cli.Timeout = time.Duration(60) * time.Second
	httpTransport := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
		MaxConnsPerHost:     100,
		IdleConnTimeout:     time.Duration(65) * time.Second,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		Dial: (&net.Dialer{
			Timeout:   time.Duration(5) * time.Second,
			KeepAlive: time.Duration(15) * time.Second,
		}).Dial,
		TLSHandshakeTimeout:   time.Duration(5) * time.Second,
		ResponseHeaderTimeout: time.Duration(60) * time.Second,
		ExpectContinueTimeout: time.Duration(5) * time.Second,
	}

	cli.Transport = httpTransport
	return cli
}
