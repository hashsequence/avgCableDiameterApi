package dataStore

import (
	"testing"
	assert "github.com/stretchr/testify/assert"
)

func TestNewDataStore(t *testing.T) {
	ds := NewDataStore()
	ds.Add(324.626)
	ds.Add(324624.626)
	ds.Add(4667.6)
	ds.Add(0.462)
	ds.Add(21324.626)
	assert.Equal(t,ds.GetAverage(), 70188.38799999999, "should be 70188.38799999999")
	a2 := []float64{324.626, 324624.626 ,4667.6, 0.462, 21324.626}
	assert.Equal(t,ds.nums, a2, "the two arrays should be equal")
	ds.Pop()
	a2 = a2[1:]
	assert.Equal(t,ds.nums, a2, "the two arrays should be equal after popping")
	ds.Pop()
	ds.Pop()
	assert.Equal(t,ds.GetAverage(), 10662.543999999983, "should be 10662.543999999983")
	ds.Pop()
	ds.Pop()
	assert.Equal(t,ds.nums, []float64{}, "the two arrays should be equal")
	ds.Pop()
}