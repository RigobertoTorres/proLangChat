package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
	//"strings"
)

var TEST = 8888

var UDPBROAD = 8875

var HOST = 8876

var CLIENT2HOST = 8877

var HOST2CLIENT = 8878

var cns []net.Conn

var hostFlag bool = true

var hostUDP *net.UDPConn

var hostAddr *net.UDPAddr

var msgchan chan string
var sentchan chan string

type Client struct {
	conn     net.Conn
	nickname string
	ch       chan string
}

func main() {
	msgchan = make(chan string)
	sentchan = make(chan string)
	makeScreen()
	go updateScreen()
	msgchan <- "Starting quickchat..."
	netInit(10)

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		curseExit()
		os.Exit(1)
	}()
	for {
		time.Sleep(3 * time.Second)
		// send test string to all connected clients
		//for c := range cns {
		//	fmt.Fprint(cns[c], "Hello fellow client!\n")
		//}
	}
}

func netInit(connectionLimit int) {
	cns = make([]net.Conn, 0, connectionLimit)

	hostFlag = hostTest()

	if !hostFlag {
		//requestConnections()
		sendClientMessage()
	}

	if hostFlag {

		go acceptConnections()
		go messageOrg()

	}
}

// broadcast message on port 8877
func sendClientMessage() {
	nameReader := bufio.NewReader(os.Stdin)
	fmt.Print("What's your username?: ")
	username, _, _ := nameReader.ReadLine()
	fmt.Print("Enter text freely to chat: ")

	for {
		reader := bufio.NewReader(os.Stdin)
		//fmt.Print("Enter text: ")
		text, _ := reader.ReadString('\n')

		//fmt.Println("Sending client message:" + text)
		conn, err := net.Dial("tcp", ":8877")
		if err != nil {
			fmt.Printf("Broadcast error: %v\n", err)
			break
		}
		var buffer bytes.Buffer
		buffer.WriteString("\n")
		buffer.WriteString(string(username))
		buffer.WriteString(": ")
		buffer.WriteString(text)
		message := buffer.String()

		conn.Write([]byte(message))
		conn.Close()
	}
}

//broadcast availability on port 8876
func requestConnections() bool {

	fmt.Println("Client: Sending connection request")
	conn, err := net.Dial("tcp", ":8876")
	if err != nil {
		fmt.Printf("Broadcast error: %v\n", err)
		return false
	}
	io.WriteString(conn, "Requesting connection")
	go readIncomingMessages(conn)
	return true
}

func hostTest() (test bool) {
	laddr := net.TCPAddr{
		IP:   nil,
		Port: TEST,
	}
	_, err := net.ListenTCP("tcp", &laddr)
	if err != nil {
		test = !requestConnections()
		if !test {
			fmt.Println("Client!!")
		}
	} else {
		fmt.Println("Host!!")
		test = true
	}

	return
}

// listen for tcp responses on 8876
func acceptConnections() {
	fmt.Printf("Listening for connections")
	ln, err := net.Listen("tcp", ":8876")
	if err != nil {
		fmt.Printf("Can't accept connections! Err: %v\n", err)
	} else {
		for {
			conn, err1 := ln.Accept()
			fmt.Println("Accepted tcp connection.")
			if err1 != nil {
				fmt.Printf("Error accepting connections: %v\n", err1)
			} else {
				fmt.Printf("Accepted tcp connection from %s\n", conn.RemoteAddr())

				cns = append(cns, conn)
				go receiveClientMessage(msgchan, conn)

			}
		}
	}
}

func receiveClientMessage(msgchan chan<- string, clConn net.Conn) {
	fmt.Println("CL2HST: Listening for clients...")
	ln, err := net.Listen("tcp", ":8877")
	if err != nil {
		fmt.Printf("Error responding to UDP broadcast: %v\n", err)
	} else {
		for {
			buff := make([]byte, 2048)
			fmt.Println("CL2HST: Accepting client messages...")
			conn, err1 := ln.Accept()
			if err1 != nil {
				fmt.Println("CL2HOST: Error in accepting client message..")
			}
			if err1 == nil {
				fmt.Println("CL2HST: Reading Client Message...")
				conn.Read(buff)
				//					var buffer bytes.Buffer
				//					buffer.WriteString("\n")
				//					buffer.WriteString(conn.LocalAddr().String())
				//					buffer.WriteString(": ")
				//					buffer.WriteString(string(buff))
				//					buffer.WriteString("\n")
				line := string(buff)
				msgchan <- line
				fmt.Printf("CL2HST: Recieved msg '%s'\n", line)

			}
		}
	}
}

func messageOrg() {
	for {
		msg := <-msgchan
		fmt.Printf("New message: %s\n", msg)
		for _, c := range cns {
			io.WriteString(c, msg)
		}
	}
}

func readIncomingMessages(c net.Conn) {
	//bufc := bufio.NewReader(c)
	for {
		buff := make([]byte, 2048)
		//line, _, err := bufc.ReadLine()
		c.Read(buff)
		//		if err != nil {
		//			fmt.Printf("TCP read error, closing connection: %v\n", err)
		//			c.Close()
		//			return
		//		}

		if buff != nil {
			msgchan <- string(buff)

		}
	}
}
