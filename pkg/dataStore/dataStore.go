package dataStore
import(
	"sync"
	"fmt"
)

type movingAverage struct {
	nums []float64
	sum float64
}

func NewMovingAverage() *movingAverage {
    return &movingAverage{
		nums: []float64{}, 
		sum: 0,
	}
}

func (this *movingAverage) pop() {
	if len(this.nums) > 0 {
		fmt.Printf("popping %v\n", this.nums[0])
		this.sum -= this.nums[0]
		this.nums[0] = 0
		this.nums = this.nums[1:]
	}
}

func (this *movingAverage) next(val float64) {
	this.sum += val
	this.nums = append(this.nums, val)
}

type DataStore struct {
	sync.RWMutex
	*movingAverage
}


func NewDataStore() *DataStore {

	return &DataStore {
		movingAverage : NewMovingAverage(),
	}
}

func (this *DataStore) Add(e float64) {
	this.Lock()
	defer func() {
		this.Unlock()
	}()
	this.next(e)
}

func (this *DataStore) Pop() {
	this.Lock()
	defer func() {
		this.Unlock()
	}()
	this.pop()
}

func (this *DataStore) GetAverage() float64 {
	this.RLock()
	defer func() {
		this.RUnlock()
	}()
	fmt.Println("sum: ",this.sum," numCount:", len(this.nums))
	return float64(this.sum) / float64(len(this.nums))
}


