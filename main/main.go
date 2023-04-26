package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/antony-raul/urlshort"
)

func main() {
	mux := defaultMux()

	// Construa o MapHandler usando o mux como fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	// Construa o YAML Handler usando o mapHandler como // fallback
	// 	yaml := `
	// - path: /urlshort
	//   url: https://github.com/gophercises/urlshort
	// - path: /urlshort-final
	//   url: https://github.com/gophercises/urlshort/tree/final
	// `
	yaml, err := urlshort.ReadYaml("../data/teste.yaml")
	if err != nil {
		log.Fatal(err)
	}
	yamlHandler, err := urlshort.YAMLHandler(yaml, mapHandler)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", yamlHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	// mux.HandleFunc("/", hello)
	return mux
}

// func hello(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintln(w, "Hello, world!")
// }
