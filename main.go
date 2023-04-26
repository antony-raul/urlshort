package main

import (
	"fmt"
	"net/http"

	"github.com/antony-raul/urlshort/handler"
	"github.com/gorilla/mux"
)

func main() {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/{path}", handler.Redirect).Methods("GET")
	muxRouter.HandleFunc("/url", handler.CadastrarUrl).Methods("POST")

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", muxRouter)
}
