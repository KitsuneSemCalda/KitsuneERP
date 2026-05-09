package main

import (
	"log"
	"net/http"

	"KitsuneSemCalda/KitsuneERP/internal/web/router"
	"KitsuneSemCalda/KitsuneERP/internal/web/router/middlewares"
)

const Addr string = ":8080"

func main() {
	r := router.New()
	r.Use(middlewares.Logging)

	log.Println("server running at http://localhost:8080")

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal(err)
	}
}
