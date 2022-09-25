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
	lineNumber := arguments[1]
	configData := readConfig()
	//id := configData[lineNumber][0]
	//hostAddress := configData[lineNumber][1]
	port := ":" + configData[lineNumber][2]
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
		destination := textParsed[1]
		messageReceived := textParsed[2]
		if len(arguments) > 2 {
			for i := 3; i < len(arguments); i++ {
				messageReceived = messageReceived + " " + textParsed[i]
			}

		}
		destinationAdress := configData[destination][1]
		unicastSend(destination, messageReceived, destinationAdress)

	}
	//go createTCPClient(ip)

	// Temporary
	time.Sleep(time.Second * 1000)

}

// Lines 40-46 & 52-54 from https://stackoverflow.com/questions/8757389/reading-a-file-line-by-line-in-go
// Stores the config data into a hashmap key is the line number value is an array with the data arr[0] = ID, arr[1] = hostaddress arr[2] = port
func readConfig() map[string][]string {
	file, err := os.Open("config.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	configData := make(map[string][]string)
	scanner := bufio.NewScanner(file)
	currentLineNum := 0
	configLine := ""
	for scanner.Scan() {
		configLine = string(scanner.Text())
		configLineParsed := parseLine(configLine)
		configData[strconv.Itoa(currentLineNum)] = configLineParsed
		currentLineNum++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return configData
}

// Parses line and stores it in an array
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

func unicastSend(destination string, message string, hostAddress string) {
	fmt.Println("Connecting to destination", destination, hostAddress)
	createTCPClient(hostAddress, message)

}

func unicastReceive(source string, message string) {

}
