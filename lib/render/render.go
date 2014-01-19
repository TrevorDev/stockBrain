package render

import (
	"net/http"
	"html/template"
	//"log"
)

func RenderView(w http.ResponseWriter, fileLocation string, templateMap map[string] string) {
	t, _ := template.ParseFiles("view/"+fileLocation)
	//log.Println("HIT!")
	//w.Write([]byte("Hello "))
	t.Execute(w, templateMap)
}

func RenderPublic(w http.ResponseWriter, r *http.Request, fileLocation string) {
	http.ServeFile(w, r, "public/"+fileLocation)
}