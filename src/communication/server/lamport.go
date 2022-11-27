package server

import (
	"SDR_Labo1/src/communication/protocol"
	"fmt"
	"strconv"
	"sync"
)

var (
	estampille      = 0                  // Estampille of the current server
	ackWaitGroup    sync.WaitGroup       // Semaphore to wait for all acks to be received
	lamportRegister *protocol.DataPacket // Last request or release sent by the current server
)

// TODO add a timeout to the ackWaitGroup
// lamportRequest sends a request to all servers and waits for all acks to be received
// This function is blocking until the critical section can be entered
func lamportRequest() {
	// Create the request with the updated estampille and save it the local register
	estampille++
	request := protocol.DataPacket{
		Type: protocol.REQ,
		Data: []string{strconv.Itoa(estampille)},
	}
	lamportRegister = &request
	fmt.Println("Creating request ", request)

	// Send the request to all servers
	for i := range serverListeners {
		ackWaitGroup.Add(1)

		// Create a goroutine to wait for the ack of server i
		go func(index int) {
			serverListeners[index].lamportRequestChan <- request // Send the request

			// As long as we don't have proof that we can enter the critical section, wait for a response
			for {
				l := <-serverListeners[index].lamportResponseChan
				receivedEstampille, err := strconv.Atoi(serverListeners[index].lamportRegister.Data[0])
				// err != nil should never happen as we are the ones who created the packet
				if err != nil {
					fmt.Println("Error parsing estampille:", err.Error())
				}
				registeredEstampille, err := strconv.Atoi(lamportRegister.Data[0])
				if err != nil {
					fmt.Println("Error parsing estampille:", err.Error())
				}

				// If the estampille of the response is greater than our request, the server i knows we want to enter the critical section
				if (l.Type == protocol.ACK || l.Type == protocol.REQ || l.Type == protocol.RES) && registeredEstampille < receivedEstampille {
					ackWaitGroup.Done()
					break
				}
			}
		}(i)
	}
	ackWaitGroup.Wait() // Wait for all servers to acknowledge the request
}

// lamportRelease sends a release to all servers
func lamportRelease() {
	estampille++
	release := protocol.DataPacket{
		Type: protocol.RES,
		Data: []string{strconv.Itoa(estampille)},
	}
	lamportRegister = &release
	fmt.Println("Creating release ", release)
	for i := range serverListeners {
		serverListeners[i].lamportRequestChan <- release
	}
}

// lamportReceive handles the reception of a request, release or ack from a server
func lamportReceive(listener *serverListener, request protocol.DataPacket) {

	// Sync the estampille of the current server with the received request
	receivedEstampille, _ := strconv.Atoi(request.Data[0])
	if receivedEstampille > estampille {
		estampille = receivedEstampille
	}
	estampille++

	switch request.Type {
	// If the request is a request we save it and send an ack
	case protocol.REQ:
		{
			listener.lamportRegister = &request
			listener.peerServer.write(protocol.DataPacket{
				Type: protocol.ACK,
				Data: []string{strconv.Itoa(estampille)},
			})
		}
	// If the request is an ack we save it if we don't have a saved request for this server
	// Then we inform the main thread that we received an ack
	case protocol.ACK:
		{
			if listener.lamportRegister == nil || listener.lamportRegister.Type != protocol.REQ {
				listener.lamportRegister = &request
			}
			listener.lamportResponseChan <- request
		}
	// If the request is a release we save it and inform the main thread that we received a release
	// only if we have a saved request (to check if we can enter the critical section now)
	case protocol.RES:
		{
			listener.lamportRegister = &request
			if lamportRegister != nil && lamportRegister.Type == protocol.REQ {
				listener.lamportResponseChan <- request
			}
		}
	}
}
