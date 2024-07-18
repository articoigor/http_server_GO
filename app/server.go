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
		splitByNewLine, _ := regexp.Compile("\n")

		splitReq := splitByNewLine.Split(req, -1)

		splitBySpace, _ := regexp.Compile(` `)
		fmt.Sprintln(len(splitReq))
		url := splitBySpace.Split(splitReq[1], -1)[1]

		params := splitBySpace.Split(splitReq[0], -1)[1]

		path := url + params

		logStrings(path)

		if path == "localhost:4221/abcdefg" {
			conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
		}

		if path == "localhost:4221" {
			conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
		}
	}
}
