package routes

import (
    "fmt"
	"encoding/json"
    "net/http"
    "log"
    "errors"
    "strconv"
    ds "github.com/hashsequence/avgCableDiameterApi/pkg/dataStore"
    utils "github.com/hashsequence/avgCableDiameterApi/pkg/utils"
  )

//GetAverageHandler handles the /cable-diameter endpoint
type GetAverageHandler struct {
    dataStore *ds.DataStore
    logger *log.Logger
    responseType string
}

func NewGetAverageHandler(dataStore *ds.DataStore, logger *log.Logger, responseType string) *GetAverageHandler {
    return &GetAverageHandler{
        dataStore : dataStore,
        logger : logger,
        responseType : func() string {
            if responseType == "" {
                return "json"
            } else {
                return responseType
            }
        }(),
    }
}

// /cable-diameter GET Api Response if json
type GetAverageHandlerResponse struct {
    Value float64
}

func (this *GetAverageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    //basic validate of request
    if r.URL.Path != "/cable-diameter" {
        this.logger.Println("404 not found.")
        http.Error(w, "404 not found.", http.StatusNotFound)
        return
    }

    if r.Method != "GET" {
        this.logger.Println("Method is not supported.")
        http.Error(w, "Method is not supported.", http.StatusNotFound)
        return
    }
    //get current average from datastore 
	currAverage := this.dataStore.GetAverage()
    this.logger.Printf("GetAverageHandler called, currentAverage: %v\n",currAverage)
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
    //depending on response type configuration, response will be plain text or json
    if this.responseType == "json" {
        resp, err := json.Marshal(GetAverageHandlerResponse{currAverage})
        if err != nil {
            this.logger.Println("Error Marshalling response")
            http.Error(w, "Error Marshalling response", http.StatusExpectationFailed)
            return
	    }
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.Write(resp)
    } else if this.responseType == "plain" {
        resp := strconv.FormatFloat(currAverage, 'f', -1, 64)
        w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
        w.Write([]byte(resp))
    }
}

type IndexHandler struct {}
//Index endpoint,  dused for testing and prints memory usage to stdout
func (this *IndexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/" {
        http.Error(w, "404 not found.", http.StatusNotFound)
        return
    }
    if r.Method != "GET" {
        http.Error(w, "Method is not supported.", http.StatusNotFound)
        return
    }
    fmt.Println("this is the index")
    utils.PrintMemUsage()
  }

  //Basic Recovery Middleware wrapper to handle panics from handler
func RecoveryMiddleWare(h http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        var err error
        defer func() {
            r := recover()
            if r != nil {
                switch t := r.(type) {
                case string:
                    err = errors.New(t)
                case error:
                    err = t
                default:
                    err = errors.New("Unknown error")
                }
                http.Error(w, err.Error(), http.StatusInternalServerError)
            }
        }()
        h.ServeHTTP(w, r)
	})
}