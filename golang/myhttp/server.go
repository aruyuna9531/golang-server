package myhttp

import (
	"fmt"
	"net/http"
)



func MainHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello!!")
}