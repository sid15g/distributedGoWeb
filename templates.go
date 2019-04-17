package main

import (
	"fmt"
	"net/http"
)

type template struct {}

func (t *template) printOnly(msg string) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(rw, msg)
	}
}