package main

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

func read(conn net.Conn) (string, error) {
	b := make([]byte, 4096)
	n, err := conn.Read(b)
	if err != nil {
		return "", err
	}
	return string(b[:n]), nil
}
func decompress(compressed []byte) (string, error) {
	b := bytes.NewBuffer(compressed)
	zlibReader, err := zlib.NewReader(b)
	if err != nil {
		return "", fmt.Errorf("Failed to create zlib reader: %w\n", err)
	}
	defer zlibReader.Close()
	var decompressed bytes.Buffer
	_, err = io.Copy(&decompressed, zlibReader)
	if err != nil {
		return "", fmt.Errorf("Failed to decompress data: %w\n", err)
	}
	return decompressed.String(), nil

}
func write(conn net.Conn, answer string) (int, error) {
	writer := bufio.NewWriter(conn)
	data := answer + "\n"
	number, err := writer.WriteString(data)
	if err == nil {
		err = writer.Flush()
	}
	return number, err
}
func main() {
	conn, err := net.Dial("tcp", "challenge01.root-me.org:52022")
	if err != nil {
		log.Fatalln("Error connecting to socket:", err)
	}

	for {
		content, err := read(conn)
		if err != nil {
			log.Fatalln("Error with reading data:", err)
		}
		fmt.Println("Received data:", content)
		if strings.Contains(content, "flag") {
			fmt.Println("Found flag")
			os.Exit(0)
		}
		b64 := strings.Split(content, "'")
		compressed, err := base64.StdEncoding.DecodeString(b64[1])
		if err != nil {
			log.Fatalln("Error decoding from base64:", err)
		}
		decompressed, err := decompress(compressed)
		if err != nil {
			log.Fatalln("Zlib panic:", err)
		}
		fmt.Println("Clear data:", decompressed)
		number, err := write(conn, decompressed)
		if err != nil {
			log.Fatalln("Error writing to socker", err)
		}
		fmt.Printf("Sended %d bytes\n", number)
	}
}
