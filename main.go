package main

import (
	//"github.com/TrevorDev/go-finance"
	"net/http"
	//"fmt"
	"github.com/gorilla/mux"
	//"bufio"
	//"os"
	//"log"
	"html/template"
	"./lib/render"
	"./lib/socket"
)

func runServer() {
	h := socket.GetHub()
	go h.Run();

	r := mux.NewRouter()
    r.HandleFunc("/", HomeHandler)
    r.HandleFunc("/ws", socket.WsHandler)
    r.HandleFunc("/public/{file:.*}", publicHandler)

    http.Handle("/", r)
    http.ListenAndServe(":3000", nil)
}

func main() {
	runServer();
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

func publicHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	render.RenderPublic(w,r,vars["file"]);
}
