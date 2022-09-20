package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide port number")
		return
	}
	if len(arguments) == 2 {
		fmt.Println("Please provide host:port.")
		return
	}

	PORT := ":" + arguments[1]
	CONNECT := arguments[2]

	// I'm trying to combine the TCPServer and TCPClient into one file, he said that they both can be in the same file
	go createTCPServer(PORT)

	go createTCPClient(CONNECT)

	// Temporary
	time.Sleep(time.Second * 10)

}

/*
	Creates TCP server using first command line argument most of the code is from

here https://www.linode.com/docs/guides/developing-udp-and-tcp-clients-and-servers-in-go/
*/
func createTCPServer(PORT string) {
	l, err := net.Listen("tcp", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	c, err := l.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		netData, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		if strings.TrimSpace(string(netData)) == "STOP" {
			fmt.Println("Exiting TCP server!")
			return
		}

		fmt.Print("-> ", string(netData))
		t := time.Now()
		myTime := t.Format(time.RFC3339) + "\n"
		c.Write([]byte(myTime))
	}

}

/*
	Creates TCP client using second command line argument most of the code is from here

https://www.linode.com/docs/guides/developing-udp-and-tcp-clients-and-servers-in-go/
*/
func createTCPClient(CONNECT string) {
	c, err := net.Dial("tcp", CONNECT)
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">> ")
		text, _ := reader.ReadString('\n')
		fmt.Fprintf(c, text+"\n")

		message, _ := bufio.NewReader(c).ReadString('\n')
		fmt.Print("->: " + message)
		if strings.TrimSpace(string(text)) == "STOP" {
			fmt.Println("TCP client exiting...")
			return
		}
	}

}

func unicastSend(destination string, message string) {

}

func unicastReceive(source string, message string) {

}
