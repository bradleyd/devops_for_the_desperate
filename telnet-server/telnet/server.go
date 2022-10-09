package telnet

import (
	"bufio"
	"errors"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"telnet-server/metrics"
	"time"
)

// TCPServer holds the structure of our TCP impl
type TCPServer struct {
	addr    string
	server  net.Listener
	metrics *metrics.MetricServer
	logger  *log.Logger
}

// New creates a new Telent server.
func New(addr string, metrics *metrics.MetricServer, logger *log.Logger) *TCPServer {
	return &TCPServer{addr: addr, metrics: metrics, logger: logger}
}

// Run starts the TCP Server.
func (t *TCPServer) Run() {
	var err error
	t.server, err = net.Listen("tcp", t.addr)
	defer t.Close()

	if err != nil {
		t.logger.Printf("Failed to create listener on port %s with error %v", t.addr, err)
		os.Exit(1)
	}

	t.logger.Printf("telnet-server listening on %s\n", t.server.Addr())

	for {
		conn, err := t.server.Accept()
		if err != nil {
			err = errors.New("could not accept connection")
			t.metrics.IncrementConnectionErrors()
			continue
		}
		if conn == nil {
			err = errors.New("could not create connection")
			t.metrics.IncrementConnectionErrors()
			continue
		}
		conn.Write([]byte(banner() + "\r\n"))
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

	// source IP
	srcIP := conn.RemoteAddr().(*net.TCPAddr).IP.String()
	t.logger.Printf("[IP=%s] New session", srcIP)

	// increment metrics
	t.metrics.IncrementConnectionsProcessed()
	t.metrics.IncrementActiveConnections()

	for {
		conn.Write([]byte(">"))

		// read message
		bytes, err := reader.ReadBytes(byte('\n'))

		if err != nil {
			if err != io.EOF {
				t.logger.Println("Failed to read data, err:", err)
			}
			// increment metrics
			t.metrics.IncrementConnectionErrors()
			t.metrics.DecrementActiveConnections()
			return
		}

		// match command from client
		cmd := strings.TrimRight(string(bytes), "\r\n")
		switch cmd {
		case "quit", "q":
			conn.Write([]byte("Good Bye!\r\n"))
			t.logger.Printf("[IP=%s] User quit session", srcIP)
			t.metrics.DecrementActiveConnections()
			return
		case "date", "d":
			const layout = "Mon Jan 2 15:04:05 -0700 MST 2006"
			s := "\x1b[44;37;1m" + time.Now().Format(layout) + "\033[0m"
			conn.Write([]byte(s + "\r\n"))
		case "yell for sysop", "y":
			conn.Write([]byte("Yelling for the SysOp\r\n"))
		case "dftd":
			conn.Write([]byte("You have unlocked God mode!\r\n"))
		case "help", "?":
			command := "Command Help:\r\n1) (q)uit -- quits\r\n2) (d)ate -- prints the current datetime\r\n3) (y)ell for sysop -- gets the sysop\r\n4) (?) help -- prints this message"
			conn.Write([]byte(command + "\r\n"))
		default:
			// just echo command back since we do not handle it
			newmessage := strings.ToUpper(cmd)
			// increment metrics
			t.metrics.IncrementUnknownCommands(cmd)

			conn.Write([]byte(newmessage + "\r\n"))

		}

		t.logger.Printf("[IP=%s] Requested command: %s", srcIP, bytes)
	}
}
