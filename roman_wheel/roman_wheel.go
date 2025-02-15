package main

import (
	"fmt"
	"io"
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

type rot13Reader struct {
	r io.Reader
}

func rot13(x byte) byte {
	capital := x >= 'A' && x <= 'Z'
	if !capital && (x < 'a' || x > 'z') {
		return x
	}

	x += 13

	if capital && x > 'Z' || !capital && x > 'z' {
		x -= 26
	}
	return x
}

func (r13 *rot13Reader) Read(b []byte) (int, error) {
	n, err := r13.r.Read(b)
	for i := 0; i <= n; i++ {
		b[i] = rot13(b[i])
	}
	return n, err
}
func Write(conn net.Conn, r *rot13Reader) (int, error) {
	buf := make([]byte, 1024)
	n, err := r.Read(buf)
	if err != nil && err != io.EOF {
		return 0, err
	}
	dataToSend := append(buf[:n], '\n')
	number, err := conn.Write(dataToSend)
	return number, err
}
func main() {
	conn, err := net.Dial("tcp", "challenge01.root-me.org:52021")
	if err != nil {
		log.Fatalln(err)
	}

	for {
		content, err := Read(conn)
		fmt.Println("Received content:", content)
		if err != nil {
			log.Fatalln("Error reading data:", err)
		}
		if strings.Contains(content, "flag") {
			fmt.Println("Found flag")
			os.Exit(0)
		}
		encoded := strings.Split(content, "'")
		s := strings.NewReader(encoded[1])
		r := rot13Reader{s}
		number, err := Write(conn, &r)
		if err != nil {
			log.Fatalln("Error writing:", err)
		}
		fmt.Printf("Sended %d bytes\n", number)
	}
}
