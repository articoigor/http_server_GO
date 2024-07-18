package main

import (
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"
)

// func logStrings(s string) {
// 	fmt.Println("*********************")
// 	fmt.Println(s)
// 	fmt.Println("*********************")
// }

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
		processRequest(req, conn)
	}
}

func processRequest(req string, conn net.Conn) bool {
	splitByNewLine, _ := regexp.Compile("\n")

	splitReq := splitByNewLine.Split(req, -1)

	splitBySpace, _ := regexp.Compile(` `)

	params := splitBySpace.Split(splitReq[0], -1)[1]

	url := strings.TrimSpace(splitBySpace.Split(splitReq[1], -1)[1])

	echoRegex, _ := regexp.Compile(`/echo/(.*)`)

	echo := echoRegex.FindString(params)

	if echo != "" {
		contentRegex, _ := regexp.Compile("/")

		content := contentRegex.Split(params, -1)[1]

		str := fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(content), content)

		conn.Write([]byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n" + str))
	}

	if url == "localhost:4221" && params == "/" {
		content := ""

		conn.Write([]byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n" + content))
	} else {
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}
}
