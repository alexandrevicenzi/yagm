package main

import (
	"fmt"
	"net/http"

	"github.com/alexandrevicenzi/yagm"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World!")
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	name, ok := yagm.Param(r, "name")

	if !ok {
		name = "Unknown"
	}

	fmt.Fprintf(w, "Hello %s!", name)
}

func main() {
	mux := yagm.New()
	mux.HandleFunc("^/$", homeHandler)
	mux.Handle("^/files", http.StripPrefix("/files", http.FileServer(http.Dir("./"))))
	mux.HandleFunc("^/(?P<name>[a-zA-Z]+)$", helloHandler)
	http.ListenAndServe(":8000", mux)
}
