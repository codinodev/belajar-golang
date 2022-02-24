package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func main() {
	router := httprouter.New()
	router.GET("/", func(wirter http.ResponseWriter, request *http.Request, params httprouter.Params) {
		fmt.Fprint(wirter, "Hello HttpRouter")
	})

	server := http.Server{
		Handler: router,
		Addr:    "localhost:3000",
	}
	server.ListenAndServe()
}
