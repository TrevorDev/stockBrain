package model

import (
	"github.com/TrevorDev/go-finance"
	"log"
	"time"
	"./../lib/database"
	"./../lib/config"
	"./../lib/socket"
	"database/sql"
	"strconv"
)

type Stock struct {
}


func (m Stock)LastTradePrice(symbol string, limit int) []float32 {
	ret := make([]float32, 0)
	db := database.GetDatabase()
	rows, err := db.Query("SELECT last_trade_price FROM stock_snapshot where stock_name = $1 ORDER BY time_stamp DESC limit $2", symbol, limit)
	if err != nil {
		log.Println(err)		
	}else{
		for rows.Next() {
			var price float32
			rows.Scan(&price)
			ret = append(ret, price)
		}
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

func (m Stock)GetAvg(symbol string, numHist int)float32 {
	db := database.GetDatabase()
	rows, err := db.Query("SELECT AVG(x.last_trade_price) FROM (SELECT last_trade_price FROM stock_snapshot where stock_name = $1 ORDER BY time_stamp DESC limit $2) x", symbol, numHist)
	if err != nil {
		log.Println(err)
	}else{
		var avg float32
		for rows.Next() {
			rows.Scan(&avg)
			return avg
		}
	}
	return 0
}

func  (m Stock)BuyRecommend(symbol string)bool {
	db := database.GetDatabase()
	rows, err := db.Query("SELECT recommend_buy FROM stock_snapshot where stock_name = $1 ORDER BY time_stamp DESC limit 1", symbol)
	if err != nil {
		log.Println(err)
	}else{
		var buy bool
		for rows.Next() {
			rows.Scan(&buy)
			return buy
		}
	}
	return false
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
	}else{
		for {
			if(stockMarketOpen()){
				stockInfo := (Stock{}).QueryStocks()
		 		for k, _ := range stockInfo {
		 			socket.SendStock(k, stockInfo[k][finance.Last_Trade_Price_Only])
		 			smartPrice := (Stock{}).GetAvg(k, 30)
		 			buy := 0;
		 			stockPriceVal, _ := strconv.ParseFloat(stockInfo[k][finance.Last_Trade_Price_Only],32)
		 			if float32(stockPriceVal) < smartPrice {
		 				socket.SendStock(k, "buy")
		 				buy = 1;
	 				}else{
	 					socket.SendStock(k, "sell")
	 				}
		 			db.Exec("INSERT INTO stock_snapshot (stock_name, last_trade_price, price_per_earning, recommend_buy) VALUES ($1, $2, $3, $4)", k, stockInfo[k][finance.Last_Trade_Price_Only], stockInfo[k][finance.Price_Per_Earning_Ratio],buy)
		 		}
		 	}
		 	time.Sleep(120 * time.Second)
		}	
	}
}