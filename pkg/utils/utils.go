package utils

import (
	"fmt"
	"runtime"
	"time"
	"os"
	"encoding/json"
    "log"
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
	if configFile == "" {
		return &Configuration {
            Address : "0.0.0.0:8080",
            ReadTimeout : 10,
	        WriteTimeout : 600,
            Static : "public",
            PollApi : "http://takehome-backend.oden.network/?metric=cable-diameter",
            TimeWindow : 60,
            ResponseType : "json",
        }
	}
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

func InitMsg(serviceName, address string) {
	fmt.Println(serviceName + " started at", ": " + address)
}


//calls method f every d interval until done channel is closed or recieves a value
func DoEvery(done <-chan struct{}, d time.Duration, f func()) {
	ticker := time.NewTicker(d)	
	for {
		select {
		case <-done:
			ticker.Stop()
			return
		case <-ticker.C:
			f()
		}
	}

}

func CreateLogger(file string) *log.Logger {
	return log.New(func() *os.File {
		if file != "" {
			f, err := createFile(file)
			if err != nil {
				return os.Stdout
			}
			return f
		}
		return os.Stdout
	}(), "",0)
}

//create open file for read-write logging
func createFile(name string) (*os.File, error) {
    file, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE, 0644)
    if err != nil {
        return nil, err
    }
    return file, err
}


//prints memory usage for testing purposes
func PrintMemUsage() {
        var m runtime.MemStats
        runtime.ReadMemStats(&m)
        // For info on each, see: https://golang.org/pkg/runtime/#MemStats
        fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
        fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
        fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
        fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
    return b / 1024 / 1024
}