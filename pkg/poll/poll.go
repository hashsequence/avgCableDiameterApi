package poll

import (
	"time"
	"log"
	"encoding/json"
	"net/http"
	"io/ioutil"
	"sync"
	ds "github.com/hashsequence/avgCableDiameterApi/pkg/dataStore"
	utils "github.com/hashsequence/avgCableDiameterApi/pkg/utils"
)
//response to store polled api json response
type Poll struct {
	sync.RWMutex
	pollApi string
	timeWindow time.Duration
	done chan struct{}
	dataStore *ds.DataStore
	logger *log.Logger
}


type PollResponse struct {
    Metric string
    Value float64 
}

func NewPoll(pollApi string, timeWindow time.Duration, dataStore *ds.DataStore, logger *log.Logger) *Poll {
	return &Poll {
		pollApi : pollApi,
		timeWindow : timeWindow,
		dataStore : dataStore,
		logger : logger, 
	}
}
//method to perform Get Requests to poll Api
func (this *Poll) CallApi() {
    resp, err := http.Get(this.pollApi)
    if err == nil {
        defer resp.Body.Close()
        body, readErr := ioutil.ReadAll(resp.Body)
	    if readErr == nil {
            bodyParsed := PollResponse{}
            json.Unmarshal(body,&bodyParsed)
			this.dataStore.Add(bodyParsed.Value)
			sum, numCount, movingAverage := this.dataStore.GetAllValues()
			this.logger.Printf("polledApi Value: %v\nsum: %v numCount: %v movingAverage: %v\n",bodyParsed.Value, sum, numCount, movingAverage)
        }
    } 
}

//closes the done channel to stop the polling api
//Since I don't need to Stop the service as part of the challenge, this method
//is only used for testing purposes
func (this* Poll) Stop() {
	this.Lock()
	defer func() {
		this.Unlock()
	}()
	if this.done != nil {
		close(this.done)
		//nill out done channel so that it can be restarted again
		this.done = nil
	}
	this.logger.Printf("Stopped Polling\n")
}

//start polling api
//instantiates the done channel, so that it may be closed
//polls the api and prints the new value to log 
//and add value to dataStore's buffer and prints sum, numCount, and movingAverage to log (default to stdout)
//after a designated time(default is one minute) has passed, will begin popping the oldest value every second,logging the popped value to log
func (this* Poll) Start() {
	this.Lock()
	defer func() {
		this.Unlock()
	}()
	//cannot start if its already starting
	if this.done != nil {
		return
	}
	this.done = make(chan struct{})
    go utils.DoEvery(this.done, time.Second, this.CallApi)
    go func() {
        <-time.After(this.timeWindow)
        utils.DoEvery(this.done, time.Second, func() {
			val, ok := this.dataStore.Pop()
			if ok {
				this.logger.Println("popped: ", val)
			}
		})
	}()
	this.logger.Printf("Started Polling\n")
}
