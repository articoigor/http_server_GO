package main

import (
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"
)

func main() {
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:4221")

	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	conn, err := l.Accept()

	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	arr := make([]byte, 128)

	bytes, err := conn.Read(arr)

	req := string(bytes)

	if err == nil {
		splitReq := strings.Split(req, "\r\n")

		fmt.Println(splitReq)

		re, _ := regexp.Compile(`(GET|POST)\s(.*)\s+`)

		details := re.FindString(splitReq[0])

		path := strings.Split(details, " ")[1]

		if path == "http://localhost:4221/abcdefg" {
			conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
		}

		if path == "http://localhost:4221" {
			conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
		}
	}
}
