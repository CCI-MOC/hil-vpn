package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Request body for a create-vpn api call.
type CreateVpnReq struct {
	Vlan uint16 `json:"vlan"`
}

// Response body for a (successful) create-vpn api call.
type CreateVpnResp struct {
	Key  string `json:"key"`
	Id   string `json:"id"`
	Port uint16 `json:"port"`
}

// Create an http.Handler implementing the REST API from the spec.
func makeHandler() http.Handler {
	r := mux.NewRouter()

	r.Methods("POST").Path("/vpns/new").
		HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			var args CreateVpnReq
			err := json.NewDecoder(req.Body).Decode(&args)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Invalid Request Body"))
				return
			}

			log.Println("create vpn request:", args)
		})

	r.Methods("DELETE").Path("/vpns/{id}").
		HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			id := mux.Vars(req)["id"]

			log.Println("delete vpn request:", id)
		})

	return r
}
