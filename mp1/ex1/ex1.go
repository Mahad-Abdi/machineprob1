package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {

	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide line number of config file")
		return
	}
	lineNumber, _ := strconv.Atoi(arguments[1])
	configLine := readConfig(lineNumber)
	configLineParsed := parseLine(configLine)
	//id := configLineParsed[0]
	//ip := configLineParsed[1]
	port := ":" + configLineParsed[2]
	println(port)

	// I'm trying to combine the TCPServer and TCPClient into one file, he said that they both can be in the same file
	go createTCPServer(port)
	println("Please provide a command in the form send destination message or STOP to stop proccess")
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">> ")
		text, _ := reader.ReadString('\n')
		textParsed := parseLine(text)
		//destination := textParsed[1]
		if len(textParsed) < 3 {
			println("Please provide a command in the form send destination message or STOP to stop proccess")
			return
		}
		// add error catching later
		destination, _ := strconv.Atoi(textParsed[1])
		messageReceived := textParsed[2]
		if len(arguments) > 2 {
			for i := 3; i < len(arguments); i++ {
				messageReceived = messageReceived + " " + textParsed[i]
			}

		}
		unicastSend(destination, messageReceived)

	}
	//go createTCPClient(ip)

	// Temporary
	time.Sleep(time.Second * 1000)

}

// Lines 40-46 & 52-54 from https://stackoverflow.com/questions/8757389/reading-a-file-line-by-line-in-go
func readConfig(line int) string {
	file, err := os.Open("config.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	currentLineNum := 0
	str := ""
	for scanner.Scan() {
		if line == currentLineNum {
			str = string(scanner.Text())
			return str
		}
		currentLineNum++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return str
}

func parseLine(line string) []string {
	return strings.Split(line, " ")
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
func createTCPClient(CONNECT string, message string) {
	c, err := net.Dial("tcp", CONNECT)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Dialed", CONNECT)
	reader := bufio.NewReader(os.Stdin)
	for {
		_, err := fmt.Fprintf(c, message+"\n")
		if err != nil {
			fmt.Println("Error with client", err)
			return
		}
		fmt.Println("Sent ", message)
		if strings.TrimSpace(string(message)) == "STOP" {
			fmt.Println("TCP client exiting...")
			return
		}
		message, _ = bufio.NewReader(c).ReadString('\n')
		fmt.Println("->: Received response from ", message)
		fmt.Print(">> ")
		message, _ = reader.ReadString('\n')

	}

}

func unicastSend(destination int, message string) {
	configLine := readConfig(destination)
	configLineParsed := parseLine(configLine)
	//id := configLineParsed[0]
	//ip := configLineParsed[1]
	ip := ":" + configLineParsed[2]
	fmt.Println("Connecting to destination", destination, ip)
	createTCPClient(ip, message)

}

func unicastReceive(source string, message string) {

}
