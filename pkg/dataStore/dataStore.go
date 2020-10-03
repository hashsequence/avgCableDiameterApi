package dataStore
import(
	"sync"
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
func (this *movingAverage) pop() (float64, bool) {
	if len(this.nums) > 0 {
		poppedVal := this.nums[0]
		this.sum -= this.nums[0]
		this.nums[0] = 0
		this.nums = this.nums[1:]
		return poppedVal, true
	}
	return 0.0, false
}

//adds a new value to the list of nums
func (this *movingAverage) next(val float64) float64 {
	this.sum += val
	this.nums = append(this.nums, val)
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
func (this *DataStore) Pop() (float64, bool) {
	this.Lock()
	defer func() {
		this.Unlock()
	}()
	val, ok := this.pop()
	return val, ok
}

//readlock and gets running average
func (this *DataStore) GetAverage() float64 {
	this.RLock()
	defer func() {
		this.RUnlock()
	}()
	return float64(this.sum) / float64(len(this.nums))
}

func (this *DataStore) GetAllValues() (sum float64, numCount int, movingAverage float64) {
	this.RLock()
	defer func() {
		this.RUnlock()
	}()
	sum = this.sum
	numCount = len(this.nums)
	movingAverage = float64(this.sum) / float64(len(this.nums))
	return 
}

