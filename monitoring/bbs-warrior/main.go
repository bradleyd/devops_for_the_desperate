package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

var (
	commands              = []string{"help", "d", "ls", "foo", "y"}
	connections           = 0
	errors                = 0
	errorToConnRatio      = 0.9
	defaultMaxConnections = 250
	maxCommandsToSend     = 10
)

func connectAndDoWork(wg *sync.WaitGroup) {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	connections++

	defer wg.Done()
	conn, err := net.Dial("tcp", telnetHost())
	if err != nil {
		log.Printf("connection error: %#v\n", err)
		errors++
		return
	}

	if rnd.Float32() >= float32(errorToConnRatio) {
		conn.Close()
		errors++
		return
	}

	for i := 0; i < rnd.Intn(maxCommandsToSend); i++ {
		fmt.Fprintf(conn, commands[rnd.Intn(len(commands))]+"\n")
	}

	timeSleep := rnd.Intn(5)
	time.Sleep(time.Second * time.Duration(timeSleep))
	fmt.Fprintf(conn, "quit"+"\n")
}

func maxTelnetCommands() int {
	val, ok := os.LookupEnv("MAX_TELNET_COMMANDS")
	if !ok {
		return maxCommandsToSend
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		return maxCommandsToSend
	}
	return i
}
func maxConnections() int {
	val, ok := os.LookupEnv("MAX_CONNECTIONS")
	if !ok {
		return defaultMaxConnections
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		return defaultMaxConnections
	}
	return i
}

func telnetHost() string {
	val, ok := os.LookupEnv("TELNET_HOST")
	if !ok {
		return "telnet-server"
	}

	return val
}

func main() {

	var wg sync.WaitGroup
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	numberOfConnections := r.Intn(maxConnections())
	maxCommandsToSend = maxTelnetCommands()

	log.Printf("Launching %d connection at telnet-server", numberOfConnections)
	log.Printf("Launching %d max commands telnet-server", maxCommandsToSend)

	for i := 0; i <= numberOfConnections; i++ {
		wg.Add(1)
		go connectAndDoWork(&wg)
	}

	wg.Wait()
	log.Printf("Test Complete: %d/%d -- errors/connection", errors, connections-1)

}
