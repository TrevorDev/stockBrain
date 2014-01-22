package controller

import (
	"net/http"
	"html/template"
	"./../lib/render"
)

type Home struct {
}

func (c Home)HomeHandler(w http.ResponseWriter, r *http.Request) {
	render.RenderView(w, "index.html", map[string] interface {} {"Title": "My title", "Body": "Hi this is my body","WsUrl": template.HTML("ws://"+r.Host+"/ws")});
}