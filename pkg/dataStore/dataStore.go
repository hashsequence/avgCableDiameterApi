package dataStore
import(
	"sync"
	"fmt"
)

//movingAverage stores the buffer of the accumulated values
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

//pops the oldest value from the list of nums
func (this *movingAverage) pop() {
	if len(this.nums) > 0 {
		fmt.Printf("popping %v\n", this.nums[0])
		this.sum -= this.nums[0]
		this.nums[0] = 0
		this.nums = this.nums[1:]
	}
}

//adds a new value to the list of nums
func (this *movingAverage) next(val float64) float64 {
	this.sum += val
	this.nums = append(this.nums, val)
	fmt.Println("sum: ",this.sum," numCount:", len(this.nums), " movingAverage: ", float64(this.sum) / float64(len(this.nums)))
	return float64(this.sum) / float64(len(this.nums))
}


//concurrent version of movingAverage and is exported for use outside the package
type DataStore struct {
	sync.RWMutex
	*movingAverage
}


func NewDataStore() *DataStore {

	return &DataStore {
		movingAverage : NewMovingAverage(),
	}
}

//locks the dataStore and adds new element to movingAverage
func (this *DataStore) Add(e float64) {
	this.Lock()
	defer func() {
		this.Unlock()
	}()
	this.next(e)
}

//locks dataStore and pops oldest value
func (this *DataStore) Pop() {
	this.Lock()
	defer func() {
		this.Unlock()
	}()
	this.pop()
}

//readlock and gets running average
func (this *DataStore) GetAverage() float64 {
	this.RLock()
	defer func() {
		this.RUnlock()
	}()
	return float64(this.sum) / float64(len(this.nums))
}


