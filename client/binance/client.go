package binance

import (
	"binance-get-order-book/models"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Client interface {
	GetBidsAndAsks() (*models.BidsAndAsks, error)
}

type client struct {
	BasicUrl   string
	HttpClient *http.Client
}

func NewBinanceClient(basicUrl string, httpClient *http.Client) Client {
	return client{
		BasicUrl:   basicUrl,
		HttpClient: httpClient,
	}
}

func (c client) GetBidsAndAsks() (*models.BidsAndAsks, error) {
	limit := 20        // then will be reduced until 15
	symbol := "BNBBTC" // set to this val temporarily
	urlPath := fmt.Sprintf(c.BasicUrl+"/api/v3/depth?symbol=%v&limit=%v", symbol, limit)

	request, err := http.NewRequest("GET", urlPath, nil)
	if err != nil {
		return nil, errors.New("new request: " + err.Error())
	}

	response, err := c.HttpClient.Do(request)
	if err != nil {
		return nil, errors.New("http client do error: " + err.Error())
	}
	defer response.Body.Close()

	var bidsAndAsks models.BidsAndAsks

	respBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.New("reading body error: " + err.Error())
	}

	err = json.Unmarshal(respBody, &bidsAndAsks)
	if err != nil {
		return nil, errors.New("json decoder error: " + err.Error())
	}

	return &bidsAndAsks, nil
}
