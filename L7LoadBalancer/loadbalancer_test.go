package L7LoadBalancer

import (
	"errors"
	"fmt"
	"slices"
	"testing"
)

//test the add

func TestAdd(t *testing.T) {
	loadBalancer := NewLoadBalancer(2)

	if err := loadBalancer.Add(""); !errors.Is(err, ErrInvalidAddress) {
		t.Error("Expecting error, got none")
	}

	if err := loadBalancer.Add("127.0.0.1:8080"); err != nil {
		t.Error(err)
	}

	if err := loadBalancer.Add("127.0.0.1:8081"); err != nil {
		t.Error(err)
	}

	if err := loadBalancer.Add("127.0.0.1:8081"); !errors.Is(err, ErrDuplicateAddress) {
		t.Error("expecting error, got none")
	}

	if count := loadBalancer.Count(); count != 2 {
		t.Errorf("Expecting 2, got %d", count)
	}
}

func TestRemove(t *testing.T) {
	loadBalancer := NewLoadBalancer(10)

	if err := loadBalancer.Add("127.0.0.1:8080"); err != nil {
		t.Error(err)
	}

	if err := loadBalancer.Add("127.0.0.1:8081"); err != nil {
		t.Error(err)
	}

	if err := loadBalancer.Remove("127"); !errors.Is(err, ErrAddressNotFound) {
		t.Error("expecting error, got none")
	}

	if err := loadBalancer.Remove("127.0.0.1:8081"); err != nil {
		t.Error(err)
	}

	if count := loadBalancer.Count(); count != 1 {
		t.Errorf("Expecting 1, got %d", count)
	}

	if _, isExist := loadBalancer.mapOfAddrs["127.0.0.1:8081"]; isExist {
		t.Error("address should not exist")
	}

	if slices.Contains(loadBalancer.listOfAddrs, "127.0.0.1:8081") {
		t.Error("address should not exist")
	}

	for _, addr := range loadBalancer.listOfAddrs {
		fmt.Println(addr)
	}
}

func TestNext(t *testing.T) {
	loadBalancer := NewLoadBalancer(10)

	if _, err := loadBalancer.Next(); !errors.Is(err, ErrNoAddressesRegistered) {
		t.Error("expecting error, got none")
	}

	if err := loadBalancer.Add("127.0.0.1:8080"); err != nil {
		t.Error(err)
	}

	if err := loadBalancer.Add("127.0.0.1:8081"); err != nil {
		t.Error(err)
	}

	if addr, _ := loadBalancer.Next(); addr != "127.0.0.1:8080" {
		t.Errorf("expecting address %s, got %s", "127.0.0.1:8080", addr)
	}

	if addr, _ := loadBalancer.Next(); addr != "127.0.0.1:8081" {
		t.Errorf("expecting address %s, got %s", "127.0.0.1:8081", addr)
	}

	if err := loadBalancer.Remove("127.0.0.1:8080"); err != nil {
		t.Error(err)
	}

	if addr, _ := loadBalancer.Next(); addr != "127.0.0.1:8081" {
		t.Errorf("expecting address %s, got %s", "127.0.0.1:8081", addr)
	}

	if addr, _ := loadBalancer.Next(); addr != "127.0.0.1:8081" {
		t.Errorf("expecting address %s, got %s", "127.0.0.1:8081", addr)
	}

}
