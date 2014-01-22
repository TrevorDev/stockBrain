package controller

import (
	"net/http"
	"github.com/gorilla/mux"
	"html/template"
	"./../lib/render"
)

type Stock struct {
}

func (c Stock)StockHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	render.RenderView(w, "stock.html", map[string] interface {} {"Symbol": vars["symbol"], "WsUrl": template.HTML("ws://"+r.Host+"/ws")});
}