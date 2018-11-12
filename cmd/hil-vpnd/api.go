package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/CCI-MOC/hil-vpn/internal/validate"
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

			if err := validate.CheckVlanNo(args.Vlan); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				return
			}

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

			vpnName := makeVpnName(id, port)
			keyText, err := privops.CreateVPN(vpnName, args.Vlan, port)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println("Error creating vpn: ", err)
				states.DeleteVpn(id)
				states.ReleasePort(port)
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
					// NOTE: that in this case we do *not* return the port
					// to the free pool, since we don't want another network
					// to possibly re-use the openvpn config we just created.
				} else {
					states.DeleteVpn(id)
					states.ReleasePort(port)
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
			idStr := mux.Vars(req)["id"]
			idSlice, err := hex.DecodeString(idStr)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			var id UniqueId
			if len(idSlice) != len(id[:]) {
				// id is the wrong length.
				w.WriteHeader(http.StatusBadRequest)
			}
			copy(id[:], idSlice)

			port, err := states.DeleteVpn(id)
			switch err {
			case nil:
			case ErrNoSuchVpn:
				w.WriteHeader(http.StatusBadRequest)
				return
			default:
				w.WriteHeader(http.StatusInternalServerError)
				log.Println("Unexpected error from DeleteVpn:", err)
				return
			}
			vpnName := makeVpnName(id, port)

			if err = privops.StopVPN(vpnName); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println("Error stopping vpn:", err)
				return
			}
			if err = privops.DeleteVPN(vpnName); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println("Error deleting vpn:", err)
				return
			}

			// OK; everything went through, so it's safe to flag the port as
			// available for re-use:
			states.ReleasePort(port)
		})

	return r
}

// format the vpn name as we will pass it to PrivOps.
func makeVpnName(id UniqueId, port uint16) string {
	return fmt.Sprintf("hil_vpn_id_%x_port_%d", id, port)
}
