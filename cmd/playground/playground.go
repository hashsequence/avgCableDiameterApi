package main

import(
	poll "github.com/hashsequence/avgCableDiameterApi/pkg/poll"
	utils "github.com/hashsequence/avgCableDiameterApi/pkg/utils"
	ds "github.com/hashsequence/avgCableDiameterApi/pkg/dataStore"
	"time"
)

func main() {
	config := utils.LoadConfig("")
	//instantiate dataStore for moving average
	dataStore := ds.NewDataStore()
	//create a central logger
	logger := utils.CreateLogger(config.File)
	poller := poll.NewPoll(config.PollApi, config.TimeWindow, dataStore, logger)
	done := make(chan struct{})
	go utils.DoEvery(done, time.Second, func() {
		logger.Println("AVERAGE: ",dataStore.GetAverage())
	})
	poller.Start()
	time.Sleep(5 * time.Second)
	poller.Stop()
	time.Sleep(15 * time.Second)
	poller.Start()
	time.Sleep(5 * time.Second)
}