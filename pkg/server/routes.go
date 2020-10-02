package server

import (
    "fmt"
	"encoding/json"
	"net/http"
    utils "github.com/hashsequence/avgCableDiameterApi/pkg/utils"
  )

type GetAverageHandler struct {
    server *Server
}

func newGetAverageHandler(server *Server) *GetAverageHandler {
    return &GetAverageHandler{server}
}

func (this *GetAverageHandler) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	currAverage := this.server.dataStore.GetAverage()
	this.server.logger.Printf("currentAverage: %v\n",currAverage)
    resp, err := json.MarshalIndent(Response{currAverage}, "", "  ")
    if err != nil {
        this.server.logger.Println("Error Marshalling response")
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
    w.Write(resp)
}

func index(w http.ResponseWriter, request *http.Request) {
    fmt.Println("this is the index")
    utils.PrintMemUsage()
  }