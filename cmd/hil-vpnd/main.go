package main

import (
	"net/http"
)

func main() {
	http.Handle("/", makeHandler())
	panic(http.ListenAndServe(":8080", nil))
}
