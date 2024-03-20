package handlers

import (
	"fmt"
	"net/http"
)

const materialDir = "./myhttp/html_material/"

const FlowControlMax = 5
const FlowControlTimeGap = 30

func ToClient(w http.ResponseWriter, code int, text string) {
	w.WriteHeader(code)
	_, err := w.Write([]byte(text))
	if err != nil {
		fmt.Printf("ToClient error: " + err.Error())
	}
}
