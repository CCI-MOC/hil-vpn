package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"

	"github.com/CCI-MOC/obmd/token"
)

var adminToken = token.Token{
	0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef,
	0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef,
}

// helper function for making authenticted requests; like client.Do except
// that it sets the appropriate authentication headers.
func doReq(client *http.Client, req *http.Request) (*http.Response, error) {
	password, err := adminToken.MarshalText()
	if err != nil {
		panic(err)
	}
	// Make sure the map isn't nil:
	if req.Header == nil {
		req.Header = http.Header{}
	}
	req.SetBasicAuth("admin", string(password))
	return client.Do(req)
}

// Helper for making authenticated post requests; like client.Post except
// that it sets the appropriate authentication headers.
func postReq(client *http.Client, urlStr, contentType string, body io.Reader) (*http.Response, error) {
	// If the body doesn't have a Close() method, give it a nop one:
	bodyCloser, ok := body.(io.ReadCloser)
	if !ok {
		bodyCloser = ioutil.NopCloser(body)
	}

	URL, err := url.Parse(urlStr)
	if err != nil {
		panic(err)
	}
	return doReq(client, &http.Request{
		Method: "POST",
		URL:    URL,
		Header: http.Header{
			"Content-Type": {contentType},
		},
		Body: bodyCloser,
	})
}

// Create a an httptest.Server, returning the PrivOps it will use.
//
// The server's VpnStates will be populated with available ports in the range
// 5000-5009.
func initTestServer(ops *MockPrivOps) *httptest.Server {
	daemon, err := newDaemon(config{
		AdminToken: adminToken,
		MinPort:    5000,
		MaxPort:    5009,
	}, ops)
	if err != nil {
		panic(err)
	}

	server := httptest.NewServer(daemon.handler)

	return server
}

// Test that we reject unathenticated API calls.
func TestNoAdminDeny(t *testing.T) {
	ops := NewMockPrivOps()
	server := initTestServer(ops)
	defer server.Close()
	client := server.Client()

	resp, err := client.Post(server.URL+"/vpns/new", "application/json", bytes.NewBufferString(`
		{
			"vlan": 400
		}
	`))
	if err != nil {
		t.Fatal("Making request:", err)
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("Unexpected status code: %d", resp.StatusCode)
	}
}

// Test basic successful vpn creation.
func TestCreate(t *testing.T) {
	ops := NewMockPrivOps()
	server := initTestServer(ops)
	defer server.Close()
	successfullyCreateVpn(t, 232, ops, server)
}

// Helper for TestCreate, also used as setup elsewhere.
func successfullyCreateVpn(t *testing.T, vlanNo uint16, ops *MockPrivOps, server *httptest.Server) {
	client := server.Client()
	vlanStr := strconv.Itoa(int(vlanNo))
	resp, err := postReq(client, server.URL+"/vpns/new", "application/json", bytes.NewBufferString(`
		{
			"vlan": `+vlanStr+`
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
	if vpn.vlanNo != vlanNo {
		t.Fatalf("Created VPN does not have the expected vlan; should be %d but is %d.",
			vlanNo, vpn.vlanNo)
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

// Test that starting up a server when a vpn already exists detects the vpn.
func TestCreateRestore(t *testing.T) {
	// We create a vpn, then shut down the server...
	ops := NewMockPrivOps()
	server := initTestServer(ops)
	successfullyCreateVpn(t, 232, ops, server)
	server.Close()

	// ...then start it back up again, and create another one.
	server = initTestServer(ops)
	defer server.Close()
	client := server.Client()
	resp, err := postReq(client, server.URL+"/vpns/new", "application/json", bytes.NewBufferString(`
		{
			"vlan": 300
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

	// now check that we have two vpns, and they have different ports.
	if len(ops.vpns) != 2 {
		t.Fatalf("There should be 2 vpns, but there are only %d.", len(ops.vpns))
	}
	vpns := []*vpnInfo{}
	for _, v := range ops.vpns {
		vpns = append(vpns, v)
	}
	if vpns[0].portNo == vpns[1].portNo {
		t.Fatalf("Both vpns got the same port number (%d).", vpns[0].portNo)
	}

}

// Test expected failures creating vpns
func TestCreateFail(t *testing.T) {
	badVlans := []uint16{0, 4095, 4096, 10000}
	ops := NewMockPrivOps()
	server := initTestServer(ops)
	defer server.Close()
	client := server.Client()
	for _, vlanId := range badVlans {
		reqBody := bytes.NewBufferString(fmt.Sprintf(`{"vlan": %d}`, vlanId))
		resp, err := postReq(client, server.URL+"/vpns/new", "application/json", reqBody)
		if err != nil {
			t.Fatal("Making request:", err)
		}
		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("Unexpected status code: %d (expected %d)",
				resp.StatusCode,
				http.StatusBadRequest)
		}

		if len(ops.vpns) != 0 {
			t.Fatalf("A VPN was created; vpns: %v", ops.vpns)
		}
	}
}

// Return the name that the api sever should have given to the privops,
// according to the response. This is an implementation detail; we only
// need to know about it for testing.
func expectedVpnName(resp CreateVpnResp) string {
	return fmt.Sprintf("hil_vpn_id_%s_port_%d", resp.Id, resp.Port)
}

// Test deleting vpns
func TestDelete(t *testing.T) {
	ops := NewMockPrivOps()
	server := initTestServer(ops)
	defer server.Close()
	client := server.Client()

	resp, err := postReq(client, server.URL+"/vpns/new", "application/json", bytes.NewBufferString(`
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

	resp, err = doReq(client, &http.Request{
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

func TestRestore(t *testing.T) {
	ops := NewMockPrivOps()
	server := initTestServer(ops)
	server.Close()
}
