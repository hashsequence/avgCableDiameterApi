package main 

import (
	server "github.com/hashsequence/avgCableDiameterApi/pkg/server"
)

//sample web server
func main() {
	s := server.NewServer(server.LoadConfig("config.json"))
	//s := server.NewServer(nil)
	s.ListenAndServe()
}