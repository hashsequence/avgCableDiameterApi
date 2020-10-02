package utils

import (
	"fmt"
	"runtime"
	"log"
	"os"
	"encoding/json"
	"time"
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

func InitMsg(config *Configuration) {
	fmt.Println("avgCableDiameter started at", ": " + config.Address)
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