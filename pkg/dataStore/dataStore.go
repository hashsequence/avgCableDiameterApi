package dataStore
import(
	"sync"
)

//will be implementing a ring buffer based queue
//because we know that there will only be around 60 since we are 
//since its a one minute moving average, we know we only need a fixed 
//buffer, with a ring buffer we reduce the amount of memory allocation,
//and the number of garbage collection cycles
//the ways a ring buffer works is that we can wrap the values over using the modulus operator %
//so if the buffer size is n
//the nth value will wrap to 0
//example:
//n % n = 0
const (
	DEFAULT_QUEUE_SIZE = 128
)

type queue struct {
	nums []float64 //keeps a buffer of all the polled values from oden api
	sum float64 // keepts track of the sum
	head int
	tail int
	size int
}

func (this *queue) resize() {
	nums := make([]float64, this.size * 2)
	//if the index of the tail is greater than the head
	//all the values are between the head and the tail
	if this.tail > this.head {
		copy(nums, this.nums[this.head:this.tail])
	//otherwise all values are before the tail and after the head and we just need to copy those over
	} else {
		numElementsCopied := copy(nums, this.nums[this.head:])
		copy(nums[numElementsCopied:], this.nums[:this.tail])
	}
	//re-index the head, tail for the new nums
	this.head = 0
	this.tail = this.size
	this.nums = nums
}

func newQueue() *queue {
	//initial size of queue w
	return &queue {
		nums : make([]float64, DEFAULT_QUEUE_SIZE),
	}
}

func (this *queue) add(e float64) float64 {
	//doubles the capacity if size reaches capacity
	if this.size == len(this.nums) {
		this.resize()
	}
	this.sum += e
	this.nums[this.tail] = e 
	this.tail = (this.tail + 1) % len(this.nums)
	this.size++ 
	return float64(this.sum) / float64(this.size)
}

func (this *queue) pop() (float64, bool) {
	if this.size <= 0 {
		return 0.0, false
	}
	h := this.nums[this.head]
	this.sum -= h
	this.nums[this.head] = 0.0
	this.head = (this.head+1) % len(this.nums)
	this.size--
	//resize down if capacity is more than 4 times the size
	if len(this.nums) > DEFAULT_QUEUE_SIZE && (this.size*4) <= len(this.nums) {
		this.resize()
	} 
	return h, true
}
//concurrent version of movingAverage and is exported for use outside the package
type DataStore struct {
	sync.RWMutex
	*queue
}


func NewDataStore() *DataStore {

	return &DataStore {
		queue : newQueue(),
	}
}

//locks the dataStore and adds new element to movingAverage
func (this *DataStore) Add(e float64) {
	this.Lock()
	defer func() {
		this.Unlock()
	}()
	this.add(e)
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
	return float64(this.sum) / float64(this.size)
}

func (this *DataStore) GetAllValues() (sum float64, numCount int, movingAverage float64) {
	this.RLock()
	defer func() {
		this.RUnlock()
	}()
	sum = this.sum
	numCount = this.size
	movingAverage = float64(this.sum) / float64(this.size)
	return 
}

