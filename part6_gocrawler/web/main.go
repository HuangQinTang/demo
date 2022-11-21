package main

import (
	"crawler/web/controller"
	"flag"
	"fmt"
	"net/http"
)

var port = flag.Int("port", 0, "the port for me to listen on")

func main() {
	flag.Parse()
	if *port == 0 {
		fmt.Println("must specify a port")
		return
	}

	http.Handle("/search", controller.CreateSearchResultHandler("web/view/template.html"))
	http.Handle("/", http.FileServer(http.Dir("web/view")))

	err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	if err != nil {
		panic(err)
	}
}
