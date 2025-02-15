package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func Read(conn net.Conn) (string, error) {
	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err != nil {
		return "", err
	}
	return string(buf[:n]), nil
}

func Write(conn net.Conn, content string) (int, error) {
	writer := bufio.NewWriter(conn)
	number, err := writer.WriteString(content)
	if err == nil {
		err = writer.Flush()
	}
	return number, err
}

func Decode(encoded string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}

func main() {
	conn, err := net.Dial("tcp", "challenge01.root-me.org:52023")
	if err != nil {
		log.Fatalln(err)
	}

	for {
		content, err := Read(conn)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println("Receiver content:", content)
		if strings.Contains(content, "flag") {
			fmt.Println("Found flag")
			os.Exit(0)
		}
		parts := strings.Split(content, "'")
		decoded, err := Decode(parts[1])
		if err != nil {
			log.Fatalln(err)
		}
		answer := string(decoded) + "\n"
		number, err := Write(conn, answer)
		if err != nil {
			log.Fatalln("Error writing:", err)
		}
		fmt.Println("Wrote %d bytes", number)
	}
	conn.Close()
}
