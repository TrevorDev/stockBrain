package controller

import (
	"net/http"
	"github.com/gorilla/mux"
	"html/template"
	"./../lib/render"
	"./../model"
)

type Stock struct {
}

func (c Stock)StockHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var rec string
	if(model.Stock{}).BuyRecommend(vars["symbol"]){
		rec = "BUY"
	}else{
		rec = "SELL"
	}
	render.RenderView(w, "stock.html", map[string] interface {} {"Recommendation": rec, "Symbol": vars["symbol"], "WsUrl": template.HTML("ws://"+r.Host+"/ws/"+vars["symbol"])});
}