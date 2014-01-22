package apiController

import (
	"net/http"
	"github.com/gorilla/mux"
	"./../../model"
	"strconv"
	"encoding/json"
)

type Stock struct {
}

func (c Stock)StockPriceHandler(w http.ResponseWriter, r *http.Request) {
	stockApi := model.Stock{}

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

	b, _ := json.Marshal(ResponseObject{Prices: stockApi.LastTradePrice(vars["symbol"], int(limit))})
	w.Header().Set("Content-Type", "application/json")
    w.Write(b)
}