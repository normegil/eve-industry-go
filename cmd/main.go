package main

import (
	"log"
	"net/http"
)

func main() {
	server := http.Server{
		Addr: ":18080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write([]byte("Hello World !"))
			if err != nil {
				log.Printf("error: %s", err.Error())
			}
		}),
	}
	log.Printf("server listening: %s", server.Addr)
	if err := server.ListenAndServe(); nil != err {
		panic(err)
	}
}
