package poll

import (
	"time"
	"log"
	"encoding/json"
	"net/http"
	"io/ioutil"
	ds "github.com/hashsequence/avgCableDiameterApi/pkg/dataStore"
	utils "github.com/hashsequence/avgCableDiameterApi/pkg/utils"
)
//response to store polled api json response
type Poll struct {
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
		done : make(chan struct{}),
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
//however behavior is not consistent since there would be a race condition
//on whether the polling will poll one or more times before channel is closed
//Since I don't need to Stop the service as part of the challenge, this method
//is only used for testing purposes
func (this* Poll) Stop() {
	close(this.done)
	this.done = make(chan struct{})
}

//start polling api
//polls the api and prints the new value to log 
//and add value to dataStore's buffer and prints sum, numCount, and movingAverage to log (default to stdout)
//after a designated time(default is one minute) has passed, will begin popping the oldest value every second,logging it
//popping and adding new values every second will maintain the size of the window for the movingAverage
func (this* Poll) Start() {
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
}
