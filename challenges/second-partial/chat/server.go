// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 254.
//!+

// Chat is a server that lets clients chat with each other.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

//!+broadcaster
type client chan<- string // an outgoing message channel

var (
	entering  = make(chan client)
	leaving   = make(chan client)
	messages  = make(chan string) // all incoming client messages
	users     = make([]*user, 0, 100)
	ircPrompt = "irc-server > "
	serverTimeZone  = "America/Mexico_City"
)

type user struct {
	username 	string
	ip 			string
	connTime	string
	isAdmin 	bool
	isKicked	bool
	ch 			chan string
}

func createUser(username, ip, connTime string) *user{
	return &user{username: username, ip: ip, connTime: connTime}
}

func deleteUser(c client){
	for i, user := range users {
		if user.ch == c {
			users = append(users[:i], users[i+1:]...)
		}
	}
}

func getAllUsers(ch chan string) {
	for _, user := range users {
		ch <- ircPrompt + user.username +" connected since "+ user.connTime
	}
}

func getUserByUsername(username string) *user{
	for _, user := range users {
		if user.username == username {
			return user
		}
	}
	return nil
}

func getUserInfo(usr *user) string{
	if usr != nil {
		return ircPrompt + "Username: " + usr.username + ", IP: " + usr.ip + " Connected since: " + usr.connTime
	}
	return ircPrompt + "User not found"
}

func messageUser(rec *user, sen string, msg string) string{
	if rec != nil {
		rec.ch <- sen + " (whispers) > " + msg
		return sen + " to " +  rec.username + " > " + msg
	}
	return ircPrompt + "user not found"
}

func getTime() string{
	_, err := time.LoadLocation(serverTimeZone)
	if err != nil {
		return ircPrompt + "Failed loading timezone, try again later. "
	}
	return time.Now().Format("2006-01-02 15:04:05")
}

func kickUser(rec, sen *user) string{
	if rec != nil {
		if rec == sen {
			return ircPrompt + "you cannot kick yourself, you should do some streching!"
		}
		rec.isKicked = true
		rec.ch <- ircPrompt + "you've been kicked from the server"
		leaving <- rec.ch
		messages <- ircPrompt + rec.username + " has been kicked by " + sen.username
		return ircPrompt + " watch your back, he may be coming for you next"
	}
	return ircPrompt + "unable to find user"
} 

func assignAdmin(rec *user) *user{
	if len(users) >= 1  {
		if users[0] != nil {
			fmt.Println("new admin: ", users[0].username)
			return users[0]
		}
		fmt.Println("new admin: ", users[2].username)
		return users[1]
	}
	return nil
}

//!-broadcaster
func broadcaster() {
	clients := make(map[client]bool) // all connected clients
	for {
		select {
		case msg := <-messages:
			// Broadcast incoming message to all
			// clients' outgoing message channels.
			for cli := range clients {
				cli <- msg
			}

		case cli := <-entering:
			clients[cli] = true

		case cli := <-leaving:
			delete(clients, cli)
			deleteUser(cli)
			close(cli)
		}
	}
}

func getUsername(conn net.Conn) string{
	buf := make([]byte, 512)
	n, _ := conn.Read(buf)
	return string(buf[:n])
}

func validateUsername(username string) string{
	for _, user := range users {
		if user.username == username {
			s := rand.NewSource(time.Now().UnixNano())
			return validateUsername("username" + strconv.Itoa(rand.New(s).Intn(10000)))
		}
	}
	return username
}

//!+handleConn
func handleConn(conn net.Conn) {
	ch := make(chan string) // outgoing client messages
	go clientWriter(conn, ch)

	username := getUsername(conn)
	username = validateUsername(username)
	user := createUser(username, conn.RemoteAddr().String(), getTime())
	user.ch = ch
	users = append(users, user)
	user.isKicked = false

	who := user.username
	ch <- ircPrompt + "IRC chat ready to use, welcome "+ who
	messages <- ircPrompt+ who + " has arrived and is ready to chat!"
	entering <- ch
	
	if len(users) <= 1 {
		ch <- ircPrompt + "you've been promoted to admin"
		user.isAdmin =  true
	}


	input := bufio.NewScanner(conn)
	for input.Scan() {
		tokens := strings.Split(input.Text()," ")
		switch tokens[0] {
			case "/users":
				getAllUsers(ch)
			case "/msg":
				if len(tokens) < 3 {
					ch <- ircPrompt + "Usage: /msg <usename> <message>"
				} else {
					rec := getUserByUsername(tokens[1])
					if rec !=nil {
						ch <- messageUser(rec, who, strings.Join(tokens[2:], " "))
					}else{
						ch <- ircPrompt + "Usage: /msg <usename> <message>"
					}
				}
			case "/time":
				ch <- ircPrompt + "Local Time: "+ serverTimeZone + getTime()
			case "/user":
				if len(tokens) <= 1 {
					ch <- ircPrompt + "Usage: /user <username>"
				} else {
					ch <- getUserInfo(getUserByUsername(tokens[1]))
				}
			case "/kick":
				if user.isAdmin {
					if len(tokens) == 1 {
						ch <- ircPrompt + "Usage: /kick <username>"
					} else {
						ch <- kickUser(getUserByUsername(tokens[1]), user)
					}
				} else {
					ch <- ircPrompt + " permision denied, you're not allowed to kick other users"
				}
			default:
				if !user.isKicked {
					messages <- who + " > " + input.Text()
				} else {
					conn.Close()
				}
		}
	}

	deleteUser(user.ch)

	// NOTE: ignoring potential errors from input.Err()
	if user.isAdmin {
		newAdmin := assignAdmin(user)
		if newAdmin != nil{
			newAdmin.isAdmin = true
			messages <- ircPrompt + newAdmin.username + " has been promoted to server admin"
		}
	}

	if !user.isKicked {
		leaving <- ch
		messages <- ircPrompt + who + " has left"
	}

	conn.Close()
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg) // NOTE: ignoring network errors
	}
}



//!+main
func main() {

	if len(os.Args) < 5 {
		log.Println("Error: Invalid number of arguments.")
		log.Println("Usage: ./client -host <address> -port <port>")
		os.Exit(1)
	}

	hostPtr := flag.String("host", "localhost", "Address of the host")
	portPtr := flag.String("port", "9000", "Port of the server")

	flag.Parse()

	listener, err := net.Listen("tcp", *hostPtr+":"+*portPtr)
	log.Println(ircPrompt + "Server succesfully lauched at " + *hostPtr+":"+*portPtr )
	if err != nil {
		log.Fatal(err)
	}

	go broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}