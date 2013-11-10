package main

import (
	"log"
	"net/http"
	"os"

	"github.com/blackfisk/pond/pond"
)

func main() {
	p := pond.NewPond()
	port := ":" + os.Getenv("PORT")

	http.Handle("/", p)

	log.Println("--> Listening in", port)
	err := http.ListenAndServe(port, nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
