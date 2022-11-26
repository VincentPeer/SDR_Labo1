package server

import (
	"SDR_Labo1/src/communication/protocol"
	"fmt"
	"strconv"
	"sync"
)

var (
	estampille      = 0
	ackWaitGroup    sync.WaitGroup
	lamportRegister *protocol.DataPacket
)

// TODO add a timeout to the ackWaitGroup
func lamportRequest() {
	estampille++
	request := protocol.DataPacket{
		Type: protocol.REQ,
		Data: []string{strconv.Itoa(estampille)},
	}
	lamportRegister = &request
	fmt.Println("Creating request ", request)
	for i := range serverListeners {
		ackWaitGroup.Add(1)
		go func(index int) {
			serverListeners[index].lamportRequestChan <- request
			for {
				l := <-serverListeners[index].lamportResponseChan
				receivedEstampille, err := strconv.Atoi(serverListeners[index].lamportRegister.Data[0])
				if err != nil {
					fmt.Println("Error parsing estampille:", err.Error())
				}
				registeredEstampille, err := strconv.Atoi(lamportRegister.Data[0])
				if err != nil {
					fmt.Println("Error parsing estampille:", err.Error())
				}
				if (l.Type == protocol.ACK || l.Type == protocol.REQ || l.Type == protocol.RES) && registeredEstampille < receivedEstampille {
					ackWaitGroup.Done()
					break
				}
			}
		}(i)
	}
	ackWaitGroup.Wait()
}

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

func (listener *serverListener) lamportReceiveRequest(request protocol.DataPacket) {

	receivedEstampille, _ := strconv.Atoi(request.Data[0])
	if receivedEstampille > estampille {
		estampille = receivedEstampille
	}
	estampille++

	switch request.Type {
	case protocol.REQ:
		{
			listener.lamportRegister = &request
			listener.peerServer.write(protocol.DataPacket{
				Type: protocol.ACK,
				Data: []string{strconv.Itoa(estampille)},
			})
		}
	case protocol.ACK:
		{
			if listener.lamportRegister == nil || listener.lamportRegister.Type != protocol.REQ {
				listener.lamportRegister = &request
			}
			listener.lamportResponseChan <- request
		}
	case protocol.RES:
		{
			listener.lamportRegister = &request
			if lamportRegister != nil && lamportRegister.Type == protocol.REQ {
				listener.lamportResponseChan <- request
			}
		}
	}
}
