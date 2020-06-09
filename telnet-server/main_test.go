package main_test

import (
	"net"
	"testing"
)

/* func buildTestServer() *server {
	return New("localhost:9999")
}
*/

// Below init function
func TestNETServer_Run(t *testing.T) {
	// Simply check that the server is up and can
	// accept connections.
	conn, err := net.Dial("tcp", ":2323")
	if err != nil {
		t.Error("could not connect to server: ", err)
	}
	defer conn.Close()
}
