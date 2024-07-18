package main

import (
	"fmt"
	"net"
	"os"
	"regexp"
)

func logStrings(s string) {
	fmt.Println("*********************")
	fmt.Println(s)
	fmt.Println("*********************")
}

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

	bytes := make([]byte, 128)

	_, err = conn.Read(bytes)

	req := string(bytes)

	if err == nil {
		re, _ := regexp.Compile("\n")

		splitReq := re.Split(req, -1)

		logStrings(splitReq[1])

		path := splitReq[1]

		if path == "http://localhost:4221/abcdefg" {
			conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
		}

		if path == "http://localhost:4221" {
			conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
		}
	}
}
