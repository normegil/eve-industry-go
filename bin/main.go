package main

import "net/http"

func main() {
	server := http.Server{
		Addr: ":18080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello World !"))
		}),
	}
	if err := server.ListenAndServe(); nil != err {
		panic(err)
	}
}
