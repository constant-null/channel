package pipeline

import "sync"

// Handler handles one
type Handler interface {
	Handle(input chan interface{}, output chan interface{}, wg *sync.WaitGroup)
}
