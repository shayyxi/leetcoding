package L7LoadBalancer

import (
	"errors"
	"sync"
)

//L7 load balancer
//http api server or library
//in memory
//add, remove, get the next
//fixed size/max size
//duplicate
//addresses of the backends, url address format
//thread safe
//round robin

//skip unhealthy backends
//health/status check
//observability/monitoring

var (
	ErrCapacityReached       = errors.New("maximum number of backends reached")
	ErrInvalidAddress        = errors.New("invalid address")
	ErrAddressNotFound       = errors.New("address not found")
	ErrNoAddressesRegistered = errors.New("no addresses registered")
	ErrDuplicateAddress      = errors.New("address already registered")
)

type LoadBalancerImpl struct {
	maxSize     int
	mapOfAddrs  map[string]int
	listOfAddrs []string
	cursor      int
	mutex       sync.Mutex
}

func NewLoadBalancer(maxSize int) *LoadBalancerImpl {
	if maxSize <= 0 {
		maxSize = 10
	}

	return &LoadBalancerImpl{
		maxSize:     maxSize,
		mapOfAddrs:  make(map[string]int),
		listOfAddrs: make([]string, 0),
		cursor:      0,
		mutex:       sync.Mutex{},
	}
}

func (lb *LoadBalancerImpl) Add(addr string) error {
	//validate the lb address, we can check the format too
	if addr == "" {
		return ErrInvalidAddress
	}

	lb.mutex.Lock()
	defer lb.mutex.Unlock()

	//check the duplicate
	if _, isExist := lb.mapOfAddrs[addr]; isExist {
		return ErrDuplicateAddress
	}

	//validate the list is not full
	if len(lb.listOfAddrs) == lb.maxSize {
		return ErrCapacityReached
	}

	//add to the list
	lb.listOfAddrs = append(lb.listOfAddrs, addr)

	//add to the map
	indexOfAddr := len(lb.listOfAddrs) - 1
	lb.mapOfAddrs[addr] = indexOfAddr

	return nil
}
func (lb *LoadBalancerImpl) Remove(addr string) error {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()

	//check if address exist
	if _, isExist := lb.mapOfAddrs[addr]; !isExist {
		return ErrAddressNotFound
	}

	//get the index of last addr
	//swap the addrTBD with last addr
	//update the list of addr by removing the last addr (deduplication)
	//update the index of last addr

	//update the cursor

	indexOfLastAddr := len(lb.listOfAddrs) - 1
	valOfLastAddr := lb.listOfAddrs[indexOfLastAddr]

	indexOfAddrTBD := lb.mapOfAddrs[addr]

	lb.listOfAddrs[indexOfAddrTBD] = lb.listOfAddrs[indexOfLastAddr]
	lb.listOfAddrs = lb.listOfAddrs[:indexOfLastAddr]

	lb.mapOfAddrs[valOfLastAddr] = indexOfAddrTBD
	delete(lb.mapOfAddrs, addr)

	if lb.cursor >= len(lb.listOfAddrs) {
		lb.cursor = 0
	}

	return nil
}

func (lb *LoadBalancerImpl) Next() (string, error) {
	//return empty, error if no addresses
	if len(lb.listOfAddrs) == 0 {
		return "", ErrNoAddressesRegistered
	}

	lb.mutex.Lock()
	defer lb.mutex.Unlock()

	addrToUse := lb.listOfAddrs[lb.cursor]

	lb.cursor = (lb.cursor + 1) % len(lb.listOfAddrs)

	return addrToUse, nil
}

func (lb *LoadBalancerImpl) Count() int {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()

	return len(lb.listOfAddrs)
}
