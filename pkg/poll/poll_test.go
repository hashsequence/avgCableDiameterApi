package poll

import (
	"testing"
	"time"
	assert "github.com/stretchr/testify/assert"
	utils "github.com/hashsequence/avgCableDiameterApi/pkg/utils"
	ds "github.com/hashsequence/avgCableDiameterApi/pkg/dataStore"
)

func TestPoller(t *testing.T) {
	dataStore := ds.NewDataStore()
	logger := utils.CreateLogger("")
	poller := NewPoll("http://takehome-backend.oden.network/?metric=cable-diameter", time.Duration(60) * time.Second, dataStore, logger)

	poller.Start()
	time.Sleep(5 * time.Second)
	v1 := poller.dataStore.GetAverage()
	time.Sleep(3 * time.Second)
	v2 := poller.dataStore.GetAverage()
	assert.NotEqual(t,v1,v2,"moving average should change")
	poller.Stop()
	time.Sleep(5 * time.Second)
	v1 = poller.dataStore.GetAverage()
	time.Sleep(5 * time.Second)
	v2 = poller.dataStore.GetAverage()
	assert.Equal(t,v1,v2,"moving average should not change after stopping")
	poller.Start()
	v1 = poller.dataStore.GetAverage()
	time.Sleep(5 * time.Second)
	v2 = poller.dataStore.GetAverage()
	assert.NotEqual(t,v1,v2,"moving average should change after restarting poller")
}
