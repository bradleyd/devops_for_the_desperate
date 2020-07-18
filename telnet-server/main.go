package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	connectionsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "telnet_server_connection_total",
		Help: "The total number of connections",
	})
	connectionErrors = promauto.NewCounter(prometheus.CounterOpts{
		Name: "telnet_server_connection_errors_total",
		Help: "The total number of errors",
	})
	unknownCommands = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "telnet_server_unknown_commands_total",
		Help: "The total number of unknown commands entered",
	}, []string{"command"})

	logger = log.New(os.Stdout, "telnet-server: ", log.LstdFlags)
)

func serveMetrics() {
	port := metricServerPort()
	http.Handle("/metrics", promhttp.Handler())
	logger.Printf("Metrics endpoint listening on %s\n", port)
	http.ListenAndServe(port, nil)
}

func metricServerPort() string {
	port, ok := os.LookupEnv("METRIC_PORT")
	if !ok {
		port = "9000"
	}
	return fmt.Sprintf(":%s", port)
}

func telnetServerPort() string {
	port, ok := os.LookupEnv("TELNET_PORT")
	if !ok {
		port = "2323"
	}
	return fmt.Sprintf(":%s", port)
}

// TCPServer holds the structure of our TCP impl
type TCPServer struct {
	addr   string
	server net.Listener
}

// Run starts the TCP Server.
func (t *TCPServer) Run() {
	var err error
	t.server, err = net.Listen("tcp", t.addr)
	defer t.Close()

	if err != nil {
		logger.Printf("Failed to create listener on port %s with error %v", t.addr, err)
		os.Exit(1)
	}

	logger.Printf("telnet-server listening on %s\n", t.server.Addr())

	for {
		conn, err := t.server.Accept()
		if err != nil {
			err = errors.New("could not accept connection")
			connectionErrors.Inc()
			continue
		}
		if conn == nil {
			err = errors.New("could not create connection")
			connectionErrors.Inc()
			continue
		}
		conn.Write([]byte(banner() + "\n"))
		go t.handleConnections(conn)
	}
}

// Close shuts down the TCP Server
func (t *TCPServer) Close() (err error) {
	return t.server.Close()
}

// handles incoming requests
func (t *TCPServer) handleConnections(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	// increment metrics
	connectionsProcessed.Inc()

	for {
		conn.Write([]byte(">"))

		// read message
		bytes, err := reader.ReadBytes(byte('\n'))

		if err != nil {
			if err != io.EOF {
				logger.Println("Failed to read data, err:", err)
			}
			// increment metrics
			connectionErrors.Inc()
			return
		}

		// match command from client
		cmd := strings.TrimRight(string(bytes), "\r\n")
		switch cmd {
		case "quit", "q":
			conn.Write([]byte("Good Bye!\n"))
			logger.Println("User quit session")
			return
		case "date", "d":
			const layout = "Mon Jan 2 15:04:05 -0700 MST 2006"
			s := "\x1b[44;37;1m" + time.Now().Format(layout) + "\033[0m"
			conn.Write([]byte(s + "\n"))
		case "yell for sysop", "y":
			conn.Write([]byte("Yelling for the SysOp\n"))
		case "help", "?":
			command := "Command Help:\n1) (q)uit -- quits\n2) (d)ate -- prints the current datetime\n3) (y)ell for sysop -- gets the sysop\n4) (?) help -- prints this message"
			conn.Write([]byte(command + "\n"))
		default:
			// just echo command back since we do not handle it
			newmessage := strings.ToUpper(cmd)
			// increment metrics
			unknownCommands.WithLabelValues(cmd).Inc()

			conn.Write([]byte(newmessage + "\n"))

		}

		logger.Printf("Request command: %s", bytes)
	}
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
	go serveMetrics()

	tcpServer := TCPServer{addr: telnetServerPort()}
	tcpServer.Run()
}
