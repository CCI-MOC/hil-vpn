package main

import (
	"net/http"
)

func main() {
	http.Handle("/", makeHandler(PrivOpsCmd{}))
	panic(http.ListenAndServe(":8080", nil))
}
