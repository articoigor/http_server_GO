package main

import (
	"fmt"
	"log"
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

	for {
		conn, err := l.Accept()

		if err != nil {
			log.Println("Error accepting: ", err.Error())
			continue
		}

		go createConnection(conn)
	}
}

func createConnection(conn net.Conn) {
	bytes := make([]byte, 128)

	_, err := conn.Read(bytes)

	req := string(bytes)

	fmt.Println(req)

	if err == nil {
		go processRequest(req, conn)
	}
}

func checkEcho(params string, conn net.Conn) {
	echoRegex, _ := regexp.Compile(`/echo/(.*)`)

	echo := echoRegex.FindString(params)

	if echo != "" {
		contentRegex, _ := regexp.Compile("/")

		content := contentRegex.Split(params, -1)[2]

		str := fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(content), content)

		conn.Write([]byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n" + str))
	}
}

func checkUserAgent(arr []string, conn net.Conn) {
	if len(arr) > 1 {
		agent := arr[1]

		str := fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(agent), agent)

		conn.Write([]byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n" + str))
	}
}

func processRequest(req string, conn net.Conn) {
	fmt.Println(req)
	reqRegex, _ := regexp.Compile("\n")

	reqComponents := reqRegex.Split(req, -1)
	fmt.Println(len(reqComponents))

	spaceSplitter, _ := regexp.Compile(` `)

	params := spaceSplitter.Split(reqComponents[0], -1)[1]

	url := strings.TrimSpace(spaceSplitter.Split(reqComponents[1], -1)[1])

	userAgent := spaceSplitter.Split(reqComponents[2], -1)

	go checkEcho(params, conn)

	go checkUserAgent(userAgent, conn)

	returnMessage := "HTTP/1.1 200 OK\r\n\r\n"

	if url != "localhost:4221" || params != "/" {
		returnMessage = "HTTP/1.1 404 Not Found\r\n\r\n"
	}

	conn.Write([]byte(returnMessage))
}
