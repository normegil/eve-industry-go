package main

import (
	"fmt"
	"github.com/normegil/evevulcan/internal/http"
	"github.com/normegil/evevulcan/ui/web"
	"log"
	stdhttp "net/http"
)

func main() {
	webFrontend, err := web.Frontend()
	if err != nil {
		panic(fmt.Errorf("load frontend: %w", err))
	}
	routes, err := http.Routes(webFrontend)
	if err != nil {
		panic(fmt.Errorf("load routes: %w", err))
	}
	server := stdhttp.Server{
		Addr:    ":18080",
		Handler: routes,
	}

	log.Printf("server listening: %s", server.Addr)
	if err := server.ListenAndServe(); nil != err {
		panic(err)
	}
}
