package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"
)

func main() {
	//mini workshop countdown server
	log.SetFlags(log.Ltime)
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		log.Println("New client is connected :  ", conn.RemoteAddr())
		if err != nil {
			log.Println(err)
			return
		}
		go countingDownHandler(conn)
	}
}
func countingDownHandler(conn net.Conn) {
	defer func() {
		io.WriteString(conn, "Your connection will be closed by server ")
		conn.Close()
	}()
	io.WriteString(conn, "Enter number : ")
	input := bufio.NewScanner(conn)
	count := Scan(input)
	if count > 20 {
		return
	}
	for {
		io.WriteString(conn, strconv.Itoa(count)+"\n")
		time.Sleep(time.Second)
		count--
		if count < 0 {
			io.WriteString(conn, "Enter number : ")
			count = Scan(input)
			if count == 0 {
				log.Println("IP Add fe is : ", conn.RemoteAddr())
				break
			}
		}
	}
}

func Scan(input *bufio.Scanner) int {
	if ok := input.Scan(); !ok {
		fmt.Println("")
		log.Println("Cannot scan value from conn. Connection is closed by client")
		return 0
	}
	count, err := strconv.Atoi(input.Text())
	if err != nil {
		log.Println("Cannot convert value from Text to int. Value", input.Text())
		return 0
	}
	return count
}
