package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

func main() {

	var (
		port   string
		logger = log.New(os.Stdout, "telnet-server: ", log.LstdFlags)
	)

	// if args does not contain port number check environment variable
	if len(os.Args) < 2 {
		portEnv := os.Getenv("PORT")
		// if args and env are empty default to 2323
		if portEnv == "" {
			port = ":2323"
		} else {
			port = fmt.Sprintf(":%s", portEnv)
		}
	} else {
		// format port number from args
		port = fmt.Sprintf(":%s", os.Args[1])
	}

	// create a tcp listener
	listener, err := net.Listen("tcp", port)
	if err != nil {
		logger.Println("Failed to create listener on port:, err:", port, err)
		os.Exit(1)
	}

	logger.Printf("Listening on %s\n", listener.Addr())

	// listen for new connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.Println("Failed to accept connection, err:", err)
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
	for {
		conn.Write([]byte(">"))
		// read message
		bytes, err := reader.ReadBytes(byte('\n'))

		if err != nil {
			if err != io.EOF {
				logger.Println("failed to read data, err:", err)
			}
			return
		}
		// match command from client
		cmd := strings.TrimRight(string(bytes), "\r\n")
		switch cmd {
		case "quit", "Quit":
			conn.Write([]byte("Good Bye!\n"))
			return
		case "date":
			const layout = "Mon Jan 2 15:04:05 -0700 MST 2006"
			s := "\x1b[44;37;1m" + time.Now().Format(layout) + "\033[0m"
			conn.Write([]byte(s + "\n"))
		case "help":
			command := "current commands:\n1) quit -- quits\n2) date -- prints the current datetime\n3) help -- prints this message"
			conn.Write([]byte(command + "\n"))
		default:
			// just echo command back since we do not handle it
			newmessage := strings.ToUpper(cmd)
			conn.Write([]byte(newmessage + "\n"))

		}

		logger.Printf("request command: %s", bytes)
	}
}

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
