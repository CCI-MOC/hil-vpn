package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

// Create a an httptest.Server, returning the PrivOps it will use.
//
// The server's VpnStates will be populated with available ports in the range
// 5000-5009.
func initTestServer() (*MockPrivOps, *httptest.Server) {
	ops := NewMockPrivOps()
	states := newStates()

	for i := 0; i < 10; i++ {
		states.FreePorts = append(states.FreePorts, uint16(5000+i))
	}

	handler := makeHandler(ops, states)
	server := httptest.NewServer(handler)

	return ops, server
}

func TestCreate(t *testing.T) {
	ops, server := initTestServer()
	defer server.Close()

	client := server.Client()
	resp, err := client.Post(server.URL+"/vpns/new", "application/json", bytes.NewBufferString(`
		{
			"vlan": 232
		}
	`))
	if err != nil {
		t.Fatal("Making request:", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Unexpected status code: %d", resp.StatusCode)
	}
	var results CreateVpnResp
	err = json.NewDecoder(resp.Body).Decode(&results)
	if err != nil {
		t.Fatal("Decoding response body:", err)
	}

	vpn, ok := ops.vpns[expectedVpnName(results)]
	if !ok {
		t.Fatalf("API request returned success, but vpn %s does not exist.", results.Id)
	}
	if vpn.vlanNo != 232 {
		t.Fatalf("Created VPN does not have the expected vlan; should be 232 but is %d.",
			vpn.vlanNo)
	}
	if vpn.key != results.Key {
		t.Fatalf("Returned key disagrees with stored key; %v vs %v", results.Key, vpn.key)
	}
	if vpn.portNo != results.Port {
		t.Fatalf("Returned port disagrees with stored port; %v vs %v", results.Port,
			vpn.portNo)
	}
	if !vpn.running {
		t.Fatal("Returned vpn was not started.")
	}
}

// Return the name that the api sever should have given to the privops,
// according to the response. This is an implementation detail; we only
// need to know about it for testing.
func expectedVpnName(resp CreateVpnResp) string {
	return fmt.Sprintf("hil_vpn_id_%s_port_%d", resp.Id, resp.Port)
}

func TestDelete(t *testing.T) {
	ops, server := initTestServer()
	defer server.Close()
	client := server.Client()

	resp, err := client.Post(server.URL+"/vpns/new", "application/json", bytes.NewBufferString(`
		{
			"vlan": 232
		}
	`))

	var results CreateVpnResp
	err = json.NewDecoder(resp.Body).Decode(&results)
	if err != nil {
		t.Fatal("Decoding response body:", err)
	}

	_, ok := ops.vpns[expectedVpnName(results)]
	if !ok {
		t.Fatal("VPN not created")
	}

	deleteUrl, err := url.Parse(server.URL + "/vpns/" + results.Id)
	if err != nil {
		panic(err)
	}

	resp, err = client.Do(&http.Request{
		Method: "DELETE",
		URL:    deleteUrl,
	})

	if err != nil {
		t.Fatal("Error making http request:", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatal("Unexpected status code:", err)
	}

	_, ok = ops.vpns[expectedVpnName(results)]
	if ok {
		t.Fatal("VPN not deleted")
	}
}
