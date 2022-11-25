package server

import (
	"fmt"
	"strconv"
)

var (
	estampille = 0
)

func syncServers(request databaseRequest) {
	fmt.Println("Syncing servers")
	estampille++
	request.payload.Data = append(request.payload.Data, strconv.Itoa(estampille))
	for _, listener := range serverListeners {
		fmt.Println("Sending request to server: ", listener)
		listener.lamportChan <- request
	}
}
