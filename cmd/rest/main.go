package main

import "net/http"

func main() {
}

func Server() http.Handler {
	mux := http.NewServeMux()
	return mux
}
