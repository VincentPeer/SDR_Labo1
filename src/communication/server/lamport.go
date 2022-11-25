package server

import (
	"SDR_Labo1/src/communication/protocol"
	"fmt"
	"strconv"
	"sync"
)

var (
	estampille   = 0
	ackWaitGroup sync.WaitGroup
)

// TODO add a timeout to the ackWaitGroup
func lamportRequest() {
	estampille++
	request := protocol.DataPacket{
		Type: protocol.REQ,
		Data: []string{strconv.Itoa(estampille)},
	}
	fmt.Println(serverListeners)
	for _, listener := range serverListeners {
		ackWaitGroup.Add(1)
		go func(listener serverListener) {
			listener.lamportRequestChan <- request
			l := <-listener.lamportResponseChan
			if l.Type == protocol.ACK {
				ackWaitGroup.Done()
			}
		}(listener)
	}
	ackWaitGroup.Wait()
}

func (listener serverListener) lamportReceiveRequest(request protocol.DataPacket) {
	receivedEstampille, _ := strconv.Atoi(request.Data[0])
	if receivedEstampille > estampille {
		estampille = receivedEstampille
	}
	estampille++

	switch request.Type {
	case protocol.REQ:
		{
			listener.peerServer.write(protocol.DataPacket{
				Type: protocol.ACK,
				Data: []string{strconv.Itoa(estampille)},
			})
		}
	case protocol.ACK:
		{
			listener.lamportResponseChan <- request
		}
	}
}
