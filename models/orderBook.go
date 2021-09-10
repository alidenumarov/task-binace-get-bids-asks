package models

type BidsAndAsks struct {
	LastUpdateId int64      `json:"lastUpdateId"`
	Bids         [][]string `json:"bids"`
	Asks         [][]string `json:"asks"`
}

type OrderVolume struct {
	Bids float64 `json:"bids"`
	Asks float64 `json:"asks"`
}
