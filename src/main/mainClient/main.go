package main

import (
	"SDR_Labo1/src/communication/client"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

// Default connection values
var (
	connHost  = "localhost"
	connPort  = []string{"11211", "11212", "11213"}
	serverIds = []int{0, 1, 2}
)

// Main function that start a new client
// 1st arg : client name, then options are available
func main() {

	isDebug := false
	var clientName string
	var serverId int

	if len(os.Args) > 1 {
		clientName = os.Args[1]

		// Check if an id was given for the server, or select a random one otherwise
		if len(os.Args) == 2 || os.Args[2] != "-I" && os.Args[2] != "--id" {
			s1 := rand.NewSource(time.Now().UnixNano())
			r1 := rand.New(s1)
			serverId = r1.Intn(len(serverIds)) + 1
			fmt.Println(serverId)
		}

		for i := 2; i < len(os.Args); i++ {
			switch os.Args[i] {
			case "-I", "--id":
				serverId, _ = strconv.Atoi(os.Args[i+1])
				i++
			case "-H", "--host":
				connHost = os.Args[i+1]
				i++
			case "-D", "--debug":
				isDebug = true
			}
		}
	} else { // if no args given, we use default values for simple server
		serverId = 4 // Default server id
	}

	conn := client.CreateConnection(clientName, connHost, connPort[serverId], isDebug) // Create a new connection
	client.StartUI(conn)                                                               // Start the user interface
	conn.Close()
}
