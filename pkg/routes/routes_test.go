package routes

import (
	"fmt"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"time"
	"testing"
	"reflect"
	"io/ioutil"
	"strconv"
	assert "github.com/stretchr/testify/assert"
	utils "github.com/hashsequence/avgCableDiameterApi/pkg/utils"
	ds "github.com/hashsequence/avgCableDiameterApi/pkg/dataStore"
	poll "github.com/hashsequence/avgCableDiameterApi/pkg/poll"
)


//test cable-diameter api if json response
func TestCableDiameterRouteJsonResponse(t *testing.T){
	//instantiate server
	dataStore := ds.NewDataStore()
	logger := utils.CreateLogger("")
	avgCableDiameterApiHandler := NewGetAverageHandler(dataStore, logger, "json")
	poller := poll.NewPoll("http://takehome-backend.oden.network/?metric=cable-diameter", time.Duration(60) * time.Second, dataStore, logger)

	ts := httptest.NewServer(avgCableDiameterApiHandler)
    defer ts.Close()

	poller.Start()
	time.Sleep(5 * time.Second)

	req, err := http.NewRequest("GET", "/cable-diameter", nil)
    if err != nil {
        t.Fatal(err)
	}
	
	rr := httptest.NewRecorder()
	//server request
	avgCableDiameterApiHandler.ServeHTTP(rr, req)

	 // Check the status code is what we expect.
	 if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusOK)
    }

	// Check the response body is what we expect.
	body, readErr := ioutil.ReadAll(rr.Body)
	if readErr == nil {
		bodyParsed := GetAverageHandlerResponse{}
		json.Unmarshal(body,&bodyParsed)
		assert.Equal(t, reflect.TypeOf(bodyParsed.Value), reflect.TypeOf(12.532156), "response value should be float64")
	}
}

//test cable-diameter api if response was plain-text
func TestCableDiameterRoutePLainResponse(t *testing.T){
	//instantiate server
	dataStore := ds.NewDataStore()
	logger := utils.CreateLogger("")
	avgCableDiameterApiHandler := NewGetAverageHandler(dataStore, logger, "plain")
	poller := poll.NewPoll("http://takehome-backend.oden.network/?metric=cable-diameter", time.Duration(60) * time.Second, dataStore, logger)

	ts := httptest.NewServer(avgCableDiameterApiHandler)
    defer ts.Close()

	poller.Start()
	time.Sleep(5 * time.Second)

	req, err := http.NewRequest("GET", "/cable-diameter", nil)
    if err != nil {
        t.Fatal(err)
	}
	
	rr := httptest.NewRecorder()
	//server request
	avgCableDiameterApiHandler.ServeHTTP(rr, req)

	 // Check the status code is what we expect.
	 if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusOK)
    }

	// Check the response body is what we expect.
	body, readErr := ioutil.ReadAll(rr.Body)
	if readErr == nil {
		f := string(body)
		s, err := strconv.ParseFloat(f, 64)
		if err != nil {
			t.Errorf("plaintext response should be a float64 when converted %v",s)
		}
		fmt.Println("plaintext response: ", s, "type: ", reflect.TypeOf(s))
		assert.Equal(t, reflect.TypeOf(s), reflect.TypeOf(35.236236), "response value should be float64")
	}
}