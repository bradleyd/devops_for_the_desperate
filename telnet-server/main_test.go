package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"testing"
	"time"
)

/* func buildTestServer() *server {
	return New("localhost:9999")
}
*/
func server(c <-chan int) {
	cmd := exec.Command("./telnet-server")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		fmt.Printf("Error executing command: %s......\n", err.Error())
	}

	//if err := cmd.Wait(); err != nil {
	//		fmt.Printf("Error waiting for command execution: %s......\n", err.Error())
	//	}
	select {
	case <-c:
		fmt.Println("Received done from channel")
		cmd.Process.Kill()
	}
}

// Below init function
func TestNETServer_Run(t *testing.T) {
	// Simply check that the server is up and can
	// accept connections.
	done := make(chan int, 1)
	go server(done)

	time.Sleep(2 * time.Second)

	conn, err := net.Dial("tcp", "127.0.0.1:2323")
	if err != nil {
		t.Error("could not connect to server: ", err)
	}
	conn.Write([]byte("q" + "\n"))
	done <- 9
	defer conn.Close()
}
