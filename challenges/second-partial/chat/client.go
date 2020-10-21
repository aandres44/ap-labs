// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 227.

// Netcat is a simple read/write client for TCP servers.
package main

import (
	"io"
	"log"
	"net"
	"os"
	"flag"
)

//!+
func main() {

	if len(os.Args) < 5 {
		log.Println("Error: Invalid number of arguments.")
		log.Println("Usage: ./client -user <username> -server <address:port>")
		os.Exit(1)
	}

	usernamePtr := flag.String("user", "User1", "The name of the user")
	serverInfoPtr := flag.String("server", "localhost:9000", "The address and the port of the server")

	flag.Parse();
	

	conn, err := net.Dial("tcp", *serverInfoPtr)
	if err != nil {
		log.Fatal(err)
	}
	done := make(chan struct{})
	go func() {
		io.Copy(os.Stdout, conn) // NOTE: ignoring errors
		log.Println("done")
		done <- struct{}{} // signal the main goroutine
	}()
	io.WriteString(conn, *usernamePtr)
	mustCopy(conn, os.Stdin)
	conn.Close()
	<-done // wait for background goroutine to finish
}

//!-

func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}
