package main

import (
	"fmt"
	"net"
	"time"
	"io"
	"bytes"
	"bufio"
	"os"
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

//var msgchan chan string

type Client struct {
	conn     net.Conn
	nickname string
	ch       chan string
}

func main() {
	fmt.Println("Starting quickchat...")
	//var message string
	
	msgchan := make(chan string)
	//addchan := make(chan Client)
	//rmchan := make(chan Client)	
	
	//msg := make(chan string)
	
	netInit(10, msgchan)
	
	go func() {
		for {
			message := <- msgchan
			fmt.Printf("%s", message)
		}
	}()
	
	for {
		time.Sleep(3 * time.Second)
		// send test string to all connected clients
		//for c := range cns {
		//	fmt.Fprint(cns[c], "Hello fellow client!\n")
		//}
	}
}

func netInit(connectionLimit int, msgchan chan string) {
	cns = make([]net.Conn, 0, connectionLimit)
	// start broadcasting for availability using Dial()
	//requestConnections()
	// start listening for requests using Listen()
	
	
	hostFlag = hostTest()
	
	if(!hostFlag) {
		requestConnections()
		sendClientMessage()
	}
	
	if(hostFlag){
		
		go acceptConnections(msgchan)
		go messageOrg(msgchan)	
		
		fmt.Print("Press Any Button")
	    var input string
	    fmt.Scanln(&input)		
		time.Sleep(3 * time.Second)			
		//go receiveClientMessage(msgchan)
	}
}

//func sendHostMessage() {
//	fmt.Println("Host: Sending UDP broadcast")
//	conn, err := net.Dial("udp", "255.255.255.255:8878")
//	if err != nil {
//		fmt.Printf("Broadcast error: %v\n", err)
//		return
//	}
//	conn.Write([]byte("Host: Hello Client!"))
//	conn.Close()
//}

//Client Receives host message to get 
//func receiveHostMessage() {
//	// buffer
//	buff := make([]byte, 2048)
//	// listen for udp
//	udpAddr := net.UDPAddr{
//		Port: 8878,
//		IP:   net.ParseIP("255.255.255.255"),
//	}
//	fmt.Println("Client: Listening for UDP...")
//	ser, err := net.ListenUDP("udp", &udpAddr)
//	if err != nil {
//		fmt.Printf("Error responding to UDP broadcast: %v\n", err)
//	} else {
//		 
//		fmt.Println("Client: Accepting UDP...")
//		hostUDP = ser
//		_, remoteAddr, err1 := ser.ReadFromUDP(buff)
//		fmt.Println("TestingTesting")
//		if err1 == nil {
//			fmt.Printf("Recieved msg '%s'\nfrom addr: %v\n", buff, remoteAddr)
//			hostAddr = remoteAddr
//		} else {
//			fmt.Println("Client: Error: %v\n", err1)
//		}	
//	}
//}

// broadcast message on port 8877
func sendClientMessage() {
	
	for {
//		fmt.Print("Enter text: ")
//	    var input string
//	    fmt.Scanln(&input)

		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter text: ")
		text, _ := reader.ReadString('\n')
		
		fmt.Println("Sending client message:" + text)
		conn, err := net.Dial("tcp", ":8877")
		if err != nil {
			fmt.Printf("Broadcast error: %v\n", err)
			break
		}
		conn.Write([]byte(text))
		conn.Close()
	}
}

 //broadcast availability on port 8876
func requestConnections() {
	
		fmt.Println("Client: Sending connection request")
		conn, err := net.Dial("tcp", ":8876")
		if err != nil {
			fmt.Printf("Broadcast error: %v\n", err)
			return
		}
		io.WriteString(conn, "Requesting connection")	
		go readIncomingMessages(conn)
		//conn.Close()
	
}

func hostTest() (test bool){
	laddr := net.TCPAddr{
		IP:   nil,
		Port: TEST,
	}
	_, err := net.ListenTCP("tcp", &laddr)
	if err != nil {
		//fmt.Printf("Can't accept connections! Err: %v\n", err)
		fmt.Println("Client!!")
		test = false
	} else {
		fmt.Println("Host!!");
		test = true	
	}
	
	return test
}

// listen for tcp responses on 8876
func acceptConnections(msgchan chan<- string) {
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
//				bufc := bufio.NewReader(conn)
//				_, err2 := io.ReadFrom(bufc)
//				if err2 == nil {
//					line, _, _ := bufc.ReadLine()
					conn.Read(buff)	
					var buffer bytes.Buffer
					buffer.WriteString("\n")
					buffer.WriteString(clConn.RemoteAddr().String())
					buffer.WriteString(": ")
					buffer.WriteString(string(buff))
					buffer.WriteString("\n")
					line := buffer.String()
					msgchan <- line							
					fmt.Printf("CL2HST: Recieved msg '%s'\n", line)
				//}
			}
		}
	}
}

func messageOrg(msgchan <-chan string){
	for {
		msg := <- msgchan
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
		
		//temp := string(line)
		if(buff != nil) {
			fmt.Print(string(buff))
		}
		//go func() {msgchan <- temp}()
	}
}
