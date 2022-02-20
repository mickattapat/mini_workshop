package main

import (
	"io"
	"log"
	"net"
	"os"
)

func main() {
	log.SetFlags(log.Ltime)
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	log.Println("Connect to server Successfuly !")
	go io.Copy(os.Stdout, conn)
	io.Copy(conn, os.Stdin)
}
