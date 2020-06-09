package main

import (
	"bufio"
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

// ascii banner
func banner() string {
	b :=
		`
____________ ___________
|  _  \  ___|_   _|  _  \
| | | | |_    | | | | | |
| | | |  _|   | | | | | |
| |/ /| |     | | | |/ /
|___/ \_|     \_/ |___/
`
	return b
}

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

func main() {

	telnetPort := telnetServerPort()

	//serve Prometheus metrics
	go serveMetrics()

	// create a tcp listener
	listener, err := net.Listen("tcp", telnetPort)
	if err != nil {
		logger.Println("Failed to create listener on port:, err:", telnetPort, err)
		os.Exit(1)
	}

	logger.Printf("telnet-server listening on %s\n", listener.Addr())

	// listen for new connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.Println("Failed to accept connection, err:", err)
			// increment metrics
			connectionErrors.Inc()
			continue
		}

		conn.Write([]byte(banner() + "\n"))

		// handle the connection
		go handleConnection(conn, logger)
	}
}

// handleConnection handles the connection
func handleConnection(conn net.Conn, logger *log.Logger) {
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
