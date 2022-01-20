package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "3333"
	CONN_TYPE = "tcp4"
)

func main() {
	host := CONN_HOST
	port := CONN_PORT
	if len(os.Args) > 1 {
		port = os.Args[1]
	}

	l, err := net.Listen(
		CONN_TYPE,
		fmt.Sprintf(":%s", port),
	)
	if err != nil {
		fmt.Printf(
			"Error creating listener: %s\n",
			err.Error(),
		)
		os.Exit(1)
	}

	defer func(host, port string) {
		fmt.Printf("Closing %v:%v.\n", host, port)
		if err := l.Close(); err != nil {
			fmt.Printf("Error while closing %s:%s because %s.\n",
				host,
				port,
				err.Error(),
			)
		}
	}(host, port)

	fmt.Printf(
		"EchoServer started on %s.\n",
		l.Addr().String(),
	)

	handleConnections(l, host, port)
}

func handleConnections(l net.Listener, host, port string) {
	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Printf(
				"Error accepting connexion: %s.\n",
				err.Error(),
			)
			return
		}
		go func(c net.Conn, host, port string) {
			handleConnection(c, host, port)
		}(c, host, port)
	}
}

func handleConnection(c net.Conn, host, port string) {
	fmt.Printf(
		"Serving %s:%s from %s.\n",
		host,
		port,
		c.RemoteAddr().String(),
	)
	connAddrStr := c.RemoteAddr().String()
	c.Write([]byte(fmt.Sprintf(
		"Connexion with endpoint %s.\n",
		connAddrStr,
	)))
	for {
		netData, err := bufio.
			NewReader(c).
			ReadString('\n')
		if err == io.EOF {
			fmt.Printf("Client %s disconnected.\n", connAddrStr)
			return

		} else if err != nil {
			fmt.Printf("Unexpected error from %s: %s\n", connAddrStr, err)
			return
		}

		newDataTrimedSpace := strings.TrimSpace(string(netData))
		if isMessageExit(newDataTrimedSpace) {
			break
		}

		fmt.Printf(
			"Echoed \"%s\" to %s.\n",
			newDataTrimedSpace,
			c.RemoteAddr(),
		)
		c.Write([]byte(newDataTrimedSpace))
		c.Write([]byte("\n"))
	}
	defer func() {
		fmt.Printf("Closing %v.\n", c.LocalAddr())
		if err := c.Close(); err != nil {
			fmt.Printf("Error while closing %s because %s.",
				c.LocalAddr(),
				err.Error(),
			)
		}
	}()
}

func isMessageExit(s string) bool {
	msg := strings.ToLower(s)
	return msg == "stop" || msg == "exit"

}
