package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"telnet-server/metrics"
	"telnet-server/telnet"
)

var (
	logger = log.New(os.Stdout, "telnet-server: ", log.LstdFlags)
)

func telnetServerPort() string {
	port, ok := os.LookupEnv("TELNET_PORT")
	if !ok {
		port = "2323"
	}
	return fmt.Sprintf(":%s", port)
}

func metricServerPort() string {
	port, ok := os.LookupEnv("METRIC_PORT")
	if !ok {
		port = "9000"
	}
	return fmt.Sprintf(":%s", port)
}

func main() {
	var info bool
	flag.BoolVar(&info, "i", false, "Print ENV")
	flag.Parse()

	if info {
		fmt.Printf("telnet port %s\nMetrics Port: %s\n", telnetServerPort(), metricServerPort())
		os.Exit(0)
	}

	//serve Prometheus metrics
	metricServer := metrics.New(metricServerPort(), logger)
	go metricServer.ListenAndServeMetrics()
	// serve Telnet
	telnetServer := telnet.New(telnetServerPort(), metricServer, logger)
	telnetServer.Run()
}
