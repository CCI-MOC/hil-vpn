package main

import (
	"encoding/json"
	"fmt"
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
func makeHandler(privops PrivOps, states *VpnStates) http.Handler {
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

			// TODO FIXME: verify that vlan is in the allowed range.

			id, port, err := states.NewVpn()
			switch err {
			case nil:
			case ErrNoFreePorts:
				w.WriteHeader(http.StatusServiceUnavailable)
				w.Write([]byte(
					"There are no free port numbers; cannot allocate " +
						"a new network."))
				return
			default:
				log.Println("error allocating new vpn: ", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			vpnName := fmt.Sprintf("hil_vpn_id_%x_port_%d", id, port)
			keyText, err := privops.CreateVPN(vpnName, args.Vlan, port)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println("Error creating vpn: ", err)
				states.DeleteVpn(id)
				return
			}

			err = privops.StartVPN(vpnName)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println("Error starting vpn: ", err)

				// try to back out the change.
				err = privops.DeleteVPN(vpnName)
				if err != nil {
					log.Println("Error deleting vpn")
				} else {
					states.DeleteVpn(id)
				}
				return
			}

			// OK, we're good -- report the info to the caller.
			w.Header().Set("Content-Type", "application/json")
			err = json.NewEncoder(w).Encode(CreateVpnResp{
				Key:  keyText,
				Id:   fmt.Sprintf("%x", id),
				Port: port,
			})
			if err != nil {
				log.Println("Error writing data to client:", err)
			}
		})

	r.Methods("DELETE").Path("/vpns/{id}").
		HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			id := mux.Vars(req)["id"]

			log.Println("delete vpn request:", id)
		})

	return r
}
