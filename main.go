package main

import (
	"net/http"
	"fmt"
	"github.com/gorilla/mux"
	"time"
	"./lib/socket"
	"./lib/database"
	"./controller"
	"./model"
	"./controller/api"
	"strconv"
)

func runServer() {
	h := socket.GetHub()
	go h.Run();

	homeCtrl := controller.Home{}
	stockCtrl := controller.Stock{}
	publicCtrl := controller.Public{}
	stockApiCtrl := apiController.Stock{}

	r := mux.NewRouter()
    r.HandleFunc("/", homeCtrl.HomeHandler)
    r.HandleFunc("/stock/{symbol}", stockCtrl.StockHandler)
    r.HandleFunc("/ws", socket.WsHandler)
    r.HandleFunc("/public/{file:.*}", publicCtrl.PublicHandler)

    r.HandleFunc("/api/stock/price/{symbol}/{limit}", stockApiCtrl.StockPriceHandler)

    http.Handle("/", r)
    http.ListenAndServe(":3004", nil)
}

func main() {
	stock := model.Stock{}
	//runtime.GOMAXPROCS(runtime.NumCPU())
	go runServer()
	database.InitDatabase()
	startTime := time.Now()
	stock.PullStocks()
	fmt.Println("stockBrain started - "+strconv.Itoa(startTime.Hour()))
}





