package main

import (
	"github.com/TrevorDev/go-finance"
	"net/http"
	"fmt"
	"github.com/gorilla/mux"
	//"bufio"
	//"os"
	"log"
	"html/template"
	"time"
	"./lib/render"
	"./lib/socket"
	"./lib/config"
	"./lib/database"
	_ "github.com/lib/pq"
	"database/sql"
	"strconv"
	"encoding/json"
)

func runServer() {
	h := socket.GetHub()
	go h.Run();

	r := mux.NewRouter()
    r.HandleFunc("/", HomeHandler)
    r.HandleFunc("/stock/{symbol}", StockHandler)
    r.HandleFunc("/ws", socket.WsHandler)
    r.HandleFunc("/public/{file:.*}", publicHandler)

    r.HandleFunc("/api/stock/price/{symbol}/{limit}", StockPriceHandler)

    http.Handle("/", r)
    http.ListenAndServe(":3000", nil)
}

func queryStocks()map[string] map[string] string {
	stockInfo, err := finance.GetStockInfo([]string{"GOOG", "FB", "YHOO", "AAPL"},[]string{finance.Last_Trade_Price_Only,finance.Price_Per_Earning_Ratio,finance.More_Info })
	if err != nil {
		fmt.Println(err)
		return queryStocks()
	}else{
		return stockInfo
	}
}

func pullStocks() {
	db, err := sql.Open("postgres", config.DatabaseUrl)
	if err != nil {
		log.Fatal(err)
	}
	/*rows, err := db.Query("SELECT last_trade_price FROM stock_snapshot where stock_name = 'GOOG'")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(rows.Columns())
	rows.Next()
	var price float32
	rows.Scan(&price)
	fmt.Println(price)*/
	for {
		if(true){
			stockInfo := queryStocks()
	 		for k, _ := range stockInfo { 
	 			/*fmt.Println(stockInfo[k][finance.Last_Trade_Price_Only])
	 			fmt.Println(stockInfo[k][finance.Price_Per_Earning_Ratio])
	 			result, err := */
	 			db.Exec("INSERT INTO stock_snapshot (stock_name, last_trade_price, price_per_earning) VALUES ($1, $2, $3)", k, stockInfo[k][finance.Last_Trade_Price_Only], stockInfo[k][finance.Price_Per_Earning_Ratio])
	 			/*fmt.Println(result)
	 			fmt.Println(err)*/
	 		}
	 	}
	 	time.Sleep(5 * time.Second)
	}
}

func stockMarketOpen(t time.Time)bool {
	if (t.Weekday() >= 1 && t.Weekday() <= 5) && ((t.Hour() > 9 && t.Hour() <= 15) || (t.Hour() == 9 && t.Minute() >= 30)) {
		return true
	}else{
		return false
	}
}

func main() {
	//runtime.GOMAXPROCS(runtime.NumCPU())
	go runServer()
	database.InitDatabase()
	startTime := time.Now()
	//fmt.Println(stockMarketOpen(startTime))
	fmt.Println("stockBrain started - "+strconv.Itoa(startTime.Hour()))
	pullStocks()
	/*fmt.Println("To stop enter q!")
	reader := bufio.NewReader(os.Stdin)
 	out, err := finance.GetStockInfo([]string{"GOOG", "MSFT"},[]string{finance.Last_Trade_Price_Only,finance.Price_Per_Earning_Ratio,finance.More_Info })
 	if(err!=nil){
 		fmt.Println(err)
 	}else{
 		fmt.Println(out)
 	}
	for input, _ := reader.ReadString('\n'); input != "q\n"; input, _ = reader.ReadString('\n') {
		fmt.Println(input)
		fmt.Println("hit3")
	}*/
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	render.RenderView(w, "index.html", map[string] interface {} {"Title": "My title", "Body": "Hi this is my body","WsUrl": template.HTML("ws://"+r.Host+"/ws")});
}

func StockHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Println(vars["symbol"]);
	render.RenderView(w, "stock.html", map[string] interface {} {"Symbol": vars["symbol"], "WsUrl": template.HTML("ws://"+r.Host+"/ws")});
}

func publicHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	render.RenderPublic(w,r,vars["file"]);
}

func StockPriceHandler(w http.ResponseWriter, r *http.Request) {

	type ResponseObject struct {
		Prices []float32
	}


	vars := mux.Vars(r)
	limit, err := strconv.ParseInt(vars["limit"],10,0)
	if(err != nil){
		b, _ := json.Marshal(ResponseObject{Prices: []float32{}})
		w.Header().Set("Content-Type", "application/json")
	    w.Write(b)
	    return
	}
	ret := make([]float32, 0)

	
	db := database.GetDatabase()
	rows, err := db.Query("SELECT last_trade_price FROM stock_snapshot where stock_name = $1 ORDER BY time_stamp DESC limit $2", vars["symbol"], limit)
	log.Println(err)
	for rows.Next() {
		var price float32
		rows.Scan(&price)
		ret = append(ret, price)
	}
	b, _ := json.Marshal(ResponseObject{Prices: ret})
	w.Header().Set("Content-Type", "application/json")
    w.Write(b)
}