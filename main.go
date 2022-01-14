package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"os"
	"strings"
	"unicode"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "3333"
	CONN_TYPE = "tcp"
)

func main() {
	host := CONN_HOST
	port := CONN_PORT
	PORT := fmt.Sprintf(":%s", port)
	l, err := net.Listen("tcp4", PORT)
	if err != nil {
		fmt.Printf("server stopped: %s", err.Error())
		os.Exit(1)
	}
	defer func(host, port string) {
		fmt.Printf("Closing ")
		if err := l.Close(); err != nil {
			fmt.Printf("Error while closing %s:%s because %s.",
				host,
				port,
				err.Error(),
			)
		}
	}(host, port)
	fmt.Println("EchoServer started!")

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Printf(
				"Error accepting connexion: %s\n",
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
		"Connexion with endpoint %s.\n\n",
		connAddrStr,
	)))
	for {
		netData, err := bufio.
			NewReader(c).
			ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}

		newDataTrimedSpace := strings.TrimSpace(string(netData))
		if newDataTrimedSpace == "STOP" {
			fmt.Printf(
				"Lost connexion from %s.\n",
				c.LocalAddr(),
			)
			break
		}

		result := fixCase(newDataTrimedSpace)
		fmt.Printf(
			"Received \"%s\" from %s and responded \"%s\".\n",
			newDataTrimedSpace,
			c.RemoteAddr(),
			string(result),
		)
		c.Write(result)
		c.Write([]byte("\n\n"))
	}
	c.Close()
}

func fixCase(s string) []byte {
	bb := bytes.NewBuffer(nil)
	bb.Grow(len(s))
	for idx, r := range s {
		var fixedRune rune = r
		if idx%2 == 0 {
			fixedRune = unicode.ToLower(r)
		} else {
			fixedRune = unicode.ToUpper(r)
		}
		fixedRune = switchPunc(fixedRune)
		b := byte(fixedRune)
		bb.WriteByte(b)
	}
	return bb.Bytes()
}

func switchPunc(r rune) rune {
	if r == '?' {
		return '!'
	} else if r == '!' {
		return '?'
	} else if r == '.' {
		return ','
	}
	return r
}
