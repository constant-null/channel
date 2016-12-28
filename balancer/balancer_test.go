package balancer

import (
	"math/rand"
	"strconv"
	"testing"
)

func TestEmptyChannelsError(t *testing.T) {
	balancer := New()

	channel, err := balancer.Get("test")

	_, correctErrorType := err.(EmptyChannelListError)
	if channel != nil || !correctErrorType {
		t.Fail()
	}

	if err.Error() != "At least one channel should be added" {
		t.Fail()
	}
}

func TestChannelAdding(t *testing.T) {
	balancer := New()

	someChan := make(chan interface{}, 2)
	balancer.Add(someChan)

	channel, err := balancer.Get("some tag")

	if channel == nil && err != nil {
		t.Fail()
	}
}

func TestChannelBalancing(t *testing.T) {
	balancer := New()

	balancer.Add(make(chan interface{}, 1))
	balancer.Add(make(chan interface{}, 1))

	ch, _ := balancer.Get("first")
	ch <- "firstMessage"

	ch, _ = balancer.Get("second")
	ch <- "secondMessage"

	fisrtChan, _ := balancer.Get("first")
	firstMessage := <-fisrtChan

	if firstMessage != "firstMessage" {
		t.Fail()
	}

	secondChan, _ := balancer.Get("second")
	secondMessage := <-secondChan

	if secondMessage != "secondMessage" {
		t.Fail()
	}

	ch, _ = balancer.Get("third")
	ch <- "thirdMessage"

	thirdMessage := <-fisrtChan

	if thirdMessage != "thirdMessage" {
		t.Fail()
	}
}

func TestCloseChannels(t *testing.T) {
	balancer := New()

	balancer.Add(make(chan interface{}, 1))
	balancer.Add(make(chan interface{}, 1))

	balancer.Close()

	_, err := balancer.Get("some tag")

	if err == nil {
		t.Fail()
	}
}

func BenchmarkBalancer(b *testing.B) {
	balancer := New()

	for i := 1; i <= 10; i++ {
		channel := make(chan interface{})
		balancer.Add(channel)
	}

	b.ResetTimer()

	for i := 1; i <= b.N; i++ {
		pid := strconv.Itoa(rand.Intn(2000))
		balancer.Get(pid)
	}

	b.ReportAllocs()
}
