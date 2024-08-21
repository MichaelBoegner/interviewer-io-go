package main

import (
	"log"
	"net/http"

	"github.com/michaelboegner/interviewer/database"
)

type apiConfig struct {
	DB *database.DB
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	mux := http.NewServeMux()

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	apiCfg := &apiConfig{}
	mux.HandleFunc("/api/users", apiCfg.handlerUsers)
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())

}
