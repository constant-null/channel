package pipeline

import (
	"sync"
	"testing"
)

type IncHandler struct {
}

func (h IncHandler) Handle(input chan interface{}, output chan interface{}, wg *sync.WaitGroup) {
	value := (<-input).(int)
	value++
	output <- value
	wg.Done()
}

func TestOneHandler(t *testing.T) {
	pipe := Pipeline{
		NewChannel: func() chan interface{} {
			return make(chan interface{}, 1)
		},
	}

	pipe.Add(IncHandler{})

	pipe.Start()

	pipe.Input() <- 0

	resultValue := <-pipe.Output()
	if resultValue != 1 {
		t.Fail()
	}
	pipe.Wait()
}

func TestTwoHandler(t *testing.T) {
	pipe := Pipeline{
		NewChannel: func() chan interface{} {
			return make(chan interface{}, 1)
		},
	}

	pipe.Add(IncHandler{})
	pipe.Add(IncHandler{})

	pipe.Start()

	pipe.Input() <- 0

	resultValue := <-pipe.Output()
	if resultValue != 2 {
		t.Fail()
	}
	pipe.Wait()
}

func TestSeveralHandlers(t *testing.T) {
	pipe := Pipeline{
		NewChannel: func() chan interface{} {
			return make(chan interface{}, 1)
		},
	}

	pipe.Add(IncHandler{})
	pipe.Add(IncHandler{})
	pipe.Add(IncHandler{})
	pipe.Add(IncHandler{})
	pipe.Add(IncHandler{})

	pipe.Start()

	pipe.Input() <- 0

	resultValue := <-pipe.Output()
	if resultValue != 5 {
		t.Fail()
	}
	pipe.Wait()
}
