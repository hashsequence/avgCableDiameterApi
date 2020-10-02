package server

import (
    "fmt"
	"encoding/json"
    "net/http"
    utils "github.com/hashsequence/avgCableDiameterApi/pkg/utils"
  )

//GetAverageHandler handles the /cable-diameter endpoint
type GetAverageHandler struct {
    server *Server
}

func newGetAverageHandler(server *Server) *GetAverageHandler {
    return &GetAverageHandler{server}
}

func (this *GetAverageHandler) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	currAverage := this.server.dataStore.GetAverage()
    this.server.logger.Printf("currentAverage: %v\n",currAverage)
    w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
    if this.server.responseType == "json" {
        resp, err := json.MarshalIndent(Response{currAverage}, "", "  ")
        if err != nil {
            this.server.logger.Println("Error Marshalling response")
	    }
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.Write(resp)
    } else {
        resp := fmt.Sprintf("%f", currAverage)    
        w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
        w.Write([]byte(resp))
    }

}

//index endpoint, used for testing and prints memory usage to stdout
func index(w http.ResponseWriter, request *http.Request) {
    fmt.Println("this is the index")
    utils.PrintMemUsage()
  }