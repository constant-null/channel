package pipeline

import (
	"sync"
	"testing"
)

type IncHandler struct {
}

func (h IncHandler) Handle(input chan interface{}, output chan interface{}, wg *sync.WaitGroup) {
	for {
		value, ok := <-input
		if !ok {
			wg.Done()
			close(output)
			return
		}
		i := value.(int)
		i++
		output <- i
	}
}

type OneTimeHandler struct {
}

func (h OneTimeHandler) Handle(input chan interface{}, output chan interface{}, wg *sync.WaitGroup) {
	value := <-input
	i := value.(int)
	i++
	output <- i
	wg.Done()
}

func TestOneHandler(t *testing.T) {
	pipe := Pipeline{
		NewChannel: func() chan interface{} {
			return make(chan interface{}, 1)
		},
	}

	pipe.Add(IncHandler{}, 1)

	pipe.Start()

	pipe.Input() <- 0

	resultValue := <-pipe.Output()
	if resultValue != 1 {
		t.Fail()
	}
	pipe.Stop()
	pipe.Wait()
}

func TestTwoHandler(t *testing.T) {
	pipe := Pipeline{
		NewChannel: func() chan interface{} {
			return make(chan interface{}, 1)
		},
	}

	pipe.Add(IncHandler{}, 1)
	pipe.Add(IncHandler{}, 1)

	pipe.Start()

	pipe.Input() <- 0

	resultValue := <-pipe.Output()
	if resultValue != 2 {
		t.Fail()
	}

	pipe.Stop()
	pipe.Wait()
}

func TestSeveralHandlers(t *testing.T) {
	pipe := Pipeline{
		NewChannel: func() chan interface{} {
			return make(chan interface{}, 1)
		},
	}

	pipe.Add(IncHandler{}, 1)
	pipe.Add(IncHandler{}, 1)
	pipe.Add(IncHandler{}, 1)
	pipe.Add(IncHandler{}, 1)
	pipe.Add(IncHandler{}, 1)

	pipe.Start()

	pipe.Input() <- 0

	resultValue := <-pipe.Output()
	if resultValue != 5 {
		t.Fail()
	}

	pipe.Stop()
	pipe.Wait()
}

func TestDifferentScaling(t *testing.T) {
	pipe := Pipeline{
		NewChannel: func() chan interface{} {
			return make(chan interface{}, 3)
		},
	}

	pipe.Add(IncHandler{}, 1)
	pipe.Add(OneTimeHandler{}, 3)

	pipe.Start()

	for i := 0; i < 3; i++ {
		pipe.Input() <- 0
	}

	for i := 0; i < 3; i++ {
		resultValue := <-pipe.Output()
		if resultValue != 2 {
			t.Fail()
		}
	}

	pipe.Stop()
	pipe.Wait()
}

func BenchmarkSeveralHandlers(b *testing.B) {
	pipe := Pipeline{
		NewChannel: func() chan interface{} {
			return make(chan interface{}, 1)
		},
	}

	pipe.Add(IncHandler{}, 1)
	pipe.Add(IncHandler{}, 1)
	pipe.Add(IncHandler{}, 1)
	pipe.Add(IncHandler{}, 1)
	pipe.Add(IncHandler{}, 1)

	pipe.Start()

	for i := 0; i < 10; i++ {
		pipe.Input() <- 0
	}

	for i := 0; i < 10; i++ {
		resultValue := <-pipe.Output()
		if resultValue != 5 {
			b.Fail()
		}
	}

	b.ReportAllocs()
	pipe.Stop()
	pipe.Wait()
}
