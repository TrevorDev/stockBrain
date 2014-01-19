package render

import (
	"net/http"
	"html/template"
	//"log"
)

func RenderView(w http.ResponseWriter, fileLocation string, templateMap map[string] interface {}) {
	t := template.Must(template.ParseFiles("view/"+fileLocation))
	t.Execute(w, templateMap)
}

func RenderPublic(w http.ResponseWriter, r *http.Request, fileLocation string) {
	http.ServeFile(w, r, "public/"+fileLocation)
}