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
	} else {
		defer conn.Close()

		createConnection(conn, "TESTE")
	}

	// theresConns := true

	// for i := 0; theresConns; i++ {

	// 	conn, err := l.Accept()

	// 	fmt.Printf("Conection count: %d", i)

	// 	defer conn.Close()

	// 	if err != nil {
	// 		fmt.Println("Error accepting connection: ", err.Error())

	// 		theresConns = false
	// 	} else {
	// 		createConnection(conn)
	// 	}
	// }
}

func createConnection(conn net.Conn, str string) {
	fmt.Println(str)
	bytes := make([]byte, 128)
	fmt.Println(len(bytes))
	_, err := conn.Read(bytes)

	req := string(bytes)

	fmt.Println(req)

	if err == nil {
		processRequest(req, conn)
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
	fmt.Println("TESTE")
	splitReq := regexp.MustCompile("\r\n").Split(req, -1)

	spaceSplitter, _ := regexp.Compile(` `)

	params := spaceSplitter.Split(splitReq[0], -1)[1]

	url := strings.TrimSpace(spaceSplitter.Split(splitReq[1], -1)[1])

	userAgent := spaceSplitter.Split(splitReq[2], -1)

	checkEcho(params, conn)

	checkUserAgent(userAgent, conn)

	returnMessage := "HTTP/1.1 200 OK\r\n\r\n"

	if url != "localhost:4221" || params != "/" {
		returnMessage = "HTTP/1.1 404 Not Found\r\n\r\n"
	}

	conn.Write([]byte(returnMessage))
}
