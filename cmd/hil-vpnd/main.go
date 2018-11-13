package main

import (
	"net/http"
)

func main() {
	// TODO: query the state of existent vpns and populate this accordingly:
	vpnStates := newStates()

	http.Handle("/", makeHandler(PrivOpsCmd{}, vpnStates))
	panic(http.ListenAndServe(":8080", nil))
}
