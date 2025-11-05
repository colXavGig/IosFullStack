package main

import (
	"ios_full_stack/web"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	web.SetRoutesToMuxiplier(mux)
	log.Panicln(http.ListenAndServe(":8080", mux))
}
