package main

import (
	"fmt"
	// "net/url"
	"net/http"
)

func MainHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello!!")
}

func main() {
	fmt.Println("Hello!")

	http.HandleFunc("/", MainHandler)
	http.ListenAndServe("0.0.0.0:8000", nil)
}