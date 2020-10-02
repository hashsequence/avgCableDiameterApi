package main 

import (
	server "github.com/hashsequence/avgCableDiameterApi/pkg/server"
)

func main() {
	s := server.NewServer(nil)
	s.ListenAndServe()
}