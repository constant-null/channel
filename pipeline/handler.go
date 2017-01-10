package pipeline

import "sync"

// Handler handles one
type Handler interface {
	Handle(input, output chan interface{}, wg *sync.WaitGroup)
}
