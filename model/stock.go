package model

import (
	"github.com/TrevorDev/go-finance"
	"log"
	"time"
	"./../lib/database"
	"./../lib/config"
	"./../lib/socket"
	"database/sql"
)

type Stock struct {
}


func (m Stock)LastTradePrice(symbol string, limit int) []float32 {
	ret := make([]float32, 0)
	db := database.GetDatabase()
	rows, err := db.Query("SELECT last_trade_price FROM stock_snapshot where stock_name = $1 ORDER BY time_stamp DESC limit $2", symbol, limit)
	if err != nil {
		log.Println(err)		
	}

	for rows.Next() {
		var price float32
		rows.Scan(&price)
		ret = append(ret, price)
	}
	return ret
}

func (m Stock)QueryStocks()map[string] map[string] string {
	stockInfo, err := finance.GetStockInfo([]string{"GOOG", "FB", "YHOO", "AAPL"},[]string{finance.Last_Trade_Price_Only,finance.Price_Per_Earning_Ratio,finance.More_Info })
	if err != nil {
		log.Println(err)
		return (Stock{}).QueryStocks()
	}else{
		return stockInfo
	}
}

func  stockMarketOpen()bool {
	t := time.Now()
	if (t.Weekday() >= 1 && t.Weekday() <= 5) && ((t.Hour() > 9 && t.Hour() <= 15) || (t.Hour() == 9 && t.Minute() >= 30)) {
		return true
	}else{
		return false
	}
}

func  (m Stock)PullStocks() {
	db, err := sql.Open("postgres", config.DatabaseUrl)
	if err != nil {
		log.Println(err)
	}
	for {
		if(stockMarketOpen()){
			stockInfo := (Stock{}).QueryStocks()
	 		for k, _ := range stockInfo {
	 			socket.SendStock(k, stockInfo[k][finance.Last_Trade_Price_Only])
	 			db.Exec("INSERT INTO stock_snapshot (stock_name, last_trade_price, price_per_earning) VALUES ($1, $2, $3)", k, stockInfo[k][finance.Last_Trade_Price_Only], stockInfo[k][finance.Price_Per_Earning_Ratio])
	 		}
	 	}
	 	time.Sleep(60 * time.Second)
	}
}