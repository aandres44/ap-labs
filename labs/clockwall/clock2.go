// Clock2 is a concurrent TCP server that periodically writes the time.
package main

import (
	"io"
	"log"
	"net"
	"time"
	"flag"
	"os"
	"fmt"
)

func handleConn(c net.Conn, location *time.Location) {
	defer c.Close()
	for {
		_, err := io.WriteString(c, fmt.Sprintf("%-16s : %s", location.String(), time.Now().In(location).Format("15:04:05\n")))
		if err != nil {
			return // e.g., client disconnected
		}
		time.Sleep(1 * time.Second)
	}
}

func main() {
	port := flag.String("port","8080","The port of the server")
	flag.Parse()
	timeZone := os.Getenv("TZ")

	location, err := time.LoadLocation(timeZone)
	if err != nil{
		log.Fatal(err)
	}

	listener, err := net.Listen("tcp", "localhost:"+*port)
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err) // e.g., connection aborted
			continue
		}
		go handleConn(conn, location) // handle connections concurrently
	}
}
