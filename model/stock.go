package model

import "time"

type Stock struct {
	StockID    uint      `json:"stockId" gorm:"primaryKey"`
	StockName  string    `json:"stockName"`
	StockPrice uint      `json:"stockPrice"`
	LastUpdate time.Time `json:"lastUpdate"`
}
