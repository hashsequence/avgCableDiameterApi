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

  //Configuration stores all the variables to build a 
  //custom server.
  type Configuration struct {
	Address string  
	ReadTimeout int64
	WriteTimeout int64
	Static string
	PollApi string
	File string
    TimeWindow int
    ResponseType string // options: "json" or "plain"
}

//Loads the configuration from a config.json
//Sample Config.json:
//{
//    "Address"        : "0.0.0.0:8080",
//    "ReadTimeout"    : 10,
//    "WriteTimeout"   : 600,
//    "Static"         : "public",
//    "pollApi"        : "http://takehome-backend.oden.network/?metric=cable-diameter",
//    "File"           : "log.txt",
//    "TimeWindow"     : 60
//}
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

func CreateFile(name string) (*os.File, error) {
    file, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE, 0644)
    if err != nil {
        return nil, err
    }
    return file, err
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
    responseType string
}

func NewServer(config *Configuration) *Server {
    if config == nil {
        config = &Configuration {
            Address : "0.0.0.0:8000",
            ReadTimeout : 10,
	        WriteTimeout : 600,
            Static : "public",
            PollApi : "http://takehome-backend.oden.network/?metric=cable-diameter",
            TimeWindow : 60,
            ResponseType : "json",
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
        log.New(func() *os.File {
            if config.File != "" {
                file, err := CreateFile(config.File)
                if err != nil {
                    return os.Stdout
                }
                return file
            }
            return os.Stdout
        }(), "",0),
        config.PollApi,
        time.Duration(config.TimeWindow) * time.Second,
        config.ResponseType,
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
