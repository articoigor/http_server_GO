package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"log"
	"net"
	"os"
	"regexp"
	"strings"
)

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
	bytes := make([]byte, 256)

	_, err := conn.Read(bytes)

	req := string(bytes)

	if err == nil {
		go processRequest(req, conn)
	}
}

func processRequest(req string, conn net.Conn) {
	reqRegex, _ := regexp.Compile("\r\n")

	reqComponents := reqRegex.Split(req, -1)

	spaceSplitter, _ := regexp.Compile(` `)

	requestDetails := spaceSplitter.Split(reqComponents[0], -1)

	method, params := requestDetails[0], requestDetails[1]

	fmt.Println(req)

	switch method {
	case "GET":
		fmt.Println("Processing GET request !")
		processGetRequest(reqComponents, params, *spaceSplitter, conn)
	case "POST":
		fmt.Println("Processing POST request !")
		processPostRequest(reqComponents, params, conn)
	}
}

func processPostRequest(components []string, params string, conn net.Conn) {
	fileNameRegex, _ := regexp.Compile(`/`)

	name := fileNameRegex.Split(params, -1)[2]

	directory := fmt.Sprintf("/tmp/data/codecrafters.io/http-server-tester/%s", name)

	err := saveFile(directory, components[5])

	if err == nil {
		conn.Write([]byte("HTTP/1.1 201 Created\r\n\r\n"))
	} else {
		fmt.Println("Error uploading file: ", err.Error())
	}
}

func processGetRequest(components []string, params string, regex regexp.Regexp, conn net.Conn) {
	url := strings.TrimSpace(regex.Split(components[1], -1)[1])

	returnMessage := "HTTP/1.1 200 OK\r\n\r\n"

	if url != "localhost:4221" || params != "/" {
		returnMessage = "HTTP/1.1 404 Not Found\r\n\r\n"
	}

	isEcho := checkEcho(components, params, &regex, conn)

	isUserAgent := checkUserAgent(params, components, &regex, conn)

	isFile := checkFile(params, conn)

	if !isEcho && !isUserAgent && !isFile {
		conn.Write([]byte(returnMessage))
	}
}

func checkEcho(components []string, params string, regex *regexp.Regexp, conn net.Conn) bool {
	echoRegex, _ := regexp.Compile(`/echo/(.*)`)

	echo := echoRegex.FindString(params)

	if echo != "" {
		bodyRegex, _ := regexp.Compile("/")
		arr := bodyRegex.Split(echo, -1)

		test := arr[len(arr)-1]

		fmt.Println(test)
		encoder, encodedBody := "", bodyRegex.Split(echo, -1)[1]

		if len(components) >= 3 {
			arr := regex.Split(components[2], -1)

			encoder = arr[len(arr)-1]

			fmt.Println(encoder)
		}

		var str string

		if encoder == "gzip" {
			content, _ := compressString(encodedBody)

			encodedBody = string(content)

			str = fmt.Sprintf("Content-Encoding: %s\r\n", encoder)
		}

		str += fmt.Sprintf("Content-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(encodedBody), encodedBody)

		conn.Write([]byte("HTTP/1.1 200 OK\r\n" + str))

		return true
	}

	return false
}

func compressString(body string) ([]byte, error) {
	var buffer bytes.Buffer
	writer := gzip.NewWriter(&buffer)

	_, err := writer.Write([]byte(body))
	if err != nil {
		return nil, err
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func checkUserAgent(param string, agents []string, regex *regexp.Regexp, conn net.Conn) bool {
	if strings.TrimSpace(param) == "/user-agent" {
		agent := regex.Split(agents[2], -1)[1]

		content := fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(agent), agent)

		conn.Write([]byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n" + content))

		return true
	}

	return false
}

func checkFile(params string, conn net.Conn) bool {
	fileRegex, _ := regexp.Compile(`/files/(.*)`)

	filePath := fileRegex.FindString(params)

	if filePath != "" {
		contentRegex, _ := regexp.Compile("/")

		content := contentRegex.Split(params, -1)[2]

		file := locateFile("/tmp/data/codecrafters.io/http-server-tester/" + content)

		if file != "" {
			str := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n%s", len(file), file)

			conn.Write([]byte(str))

			return true
		}
	}

	return false
}

func locateFile(directory string) string {
	item, err := os.ReadFile(directory)

	if err != nil {
		return ""
	}

	return string(item)
}

func saveFile(directory, content string) error {
	return os.WriteFile(directory, []byte(strings.ReplaceAll(content, "\x00", "")), os.ModeAppend)
}
