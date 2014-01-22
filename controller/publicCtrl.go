package controller

import (
	"net/http"
	"github.com/gorilla/mux"
	"./../lib/render"
)

type Public struct {
}

func (c Public)PublicHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	render.RenderPublic(w,r,vars["file"]);
}