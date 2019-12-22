package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "welcome!")
		})
	http.HandleFunc("/about",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "about ...")
		})

	http.ListenAndServe(":8080", nil)
}
