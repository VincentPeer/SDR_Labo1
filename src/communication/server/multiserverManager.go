package server

type serverConfig struct {
	id   int    `json:"id"`
	port string `json:"port"`
	host string `json:"host"`
}

type networkConfig struct {
	nbServers int            `json:"nbServers"`
	servers   []serverConfig `json:"servers"`
}
