package server

import (
    "fmt"
	"encoding/json"
	"net/http"
    "errors"
    "log"
    "os"
    "time"
    "io/ioutil"
    utils "github.com/hashsequence/avgCableDiameterApi/pkg/utils"
    ds "github.com/hashsequence/avgCableDiameterApi/pkg/dataStore"
  )

  type Configuration struct {
	Address      string  
	ReadTimeout  int64
	WriteTimeout int64
	Static       string
	PollApi string
	File *os.File
	TimeWindow time.Duration
}


func LoadConfig(configFile string) *Configuration {
	var config Configuration
	file, err := os.Open(configFile)
	if err != nil {
		log.Fatalln("Cannot open config file", err)
	}
	decoder := json.NewDecoder(file)
	config = Configuration{}
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatalln("Cannot get configuration from file", err)
	}
	return &config
}

func InitMsg(config *Configuration) {
	fmt.Println("avgCableDiameter started at", ": " + config.Address)
}

type PollResponse struct {
    Metric string
    Value float64 
}

type Response struct {
    Value float64
}

type Server struct {
    *http.Server
    mux *http.ServeMux
    dataStore *ds.DataStore
    logger *log.Logger
    pollApi string
    timeWindow time.Duration
}

func NewServer(config *Configuration) *Server {
    if config == nil {
        config = &Configuration {
            Address : "0.0.0.0:8000",
            ReadTimeout : 10,
	        WriteTimeout : 600,
            Static : "public",
            PollApi : "http://takehome-backend.oden.network/?metric=cable-diameter",
            File : os.Stdout,
            TimeWindow : time.Minute,
        }
    }
    mux := http.NewServeMux()
    return &Server{
        &http.Server{
            Addr:           config.Address,
            Handler : mux,
            ReadTimeout:    time.Duration(config.ReadTimeout * int64(time.Second)),
            WriteTimeout:   time.Duration(config.WriteTimeout * int64(time.Second)),
            MaxHeaderBytes: 1 << 20,
        },
        mux,
        ds.NewDataStore(),
        log.New(config.File, "",0),
        config.PollApi,
        config.TimeWindow,
    }

}

func (this *Server) poll() {
    resp, err := http.Get(this.pollApi)
    if err == nil {
        defer resp.Body.Close()
        body, readErr := ioutil.ReadAll(resp.Body)
	    if readErr == nil {
            bodyParsed := PollResponse{}
            json.Unmarshal(body,&bodyParsed)
            this.dataStore.Add(bodyParsed.Value)
            this.logger.Printf("polledApi Value: %v\n",bodyParsed.Value)
        }
    } 
}

func (this *Server) Routes() {
    c := newGetAverageHandler(this)
    this.mux.HandleFunc("/", index)
    this.mux.Handle("/cable-diameter", this.recoveryMiddleWare(c))
}

func (this *Server) ListenAndServe() {
    done := make(chan struct{})
    defer func() {
        close(done)
    }()
    this.Routes()
    go utils.DoEvery(done, time.Second, this.poll)
    go func() {
        <-time.After(this.timeWindow)
        utils.DoEvery(done, time.Second, this.dataStore.Pop)
    }()
    this.Server.ListenAndServe()
    
}

func (this *Server) recoveryMiddleWare(h http.Handler) http.Handler {
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
