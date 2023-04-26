package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/{path}", Redirect).Methods("GET")
	muxRouter.HandleFunc("/url", CadastrarUrl).Methods("POST")

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", muxRouter)
}
