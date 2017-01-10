package pipeline

import "sync"

// Pipeline passes input data to handlers
type Pipeline struct {
	// NewChannel optionally specifies a function to generate
	// a new channel used in handlers communications
	NewChannel func() chan interface{}
	wg         sync.WaitGroup
	input      chan interface{}
	output     chan interface{}
	handlers   []Handler
	scale      []int
}

// Input returns channel used to send data into pipeline
func (p *Pipeline) Input() chan interface{} {
	if p.input == nil {
		p.input = p.NewChannel()
	}

	return p.input
}

// Output returns channel used to get outgoing data from pipeline
func (p *Pipeline) Output() chan interface{} {
	if p.output == nil {
		p.output = p.NewChannel()
	}

	return p.output
}

// Add adds handler to the end of pipeline
func (p *Pipeline) Add(handler Handler, scale int) {
	p.scale = append(p.scale, scale)
	p.handlers = append(p.handlers, handler)
}

// Start starting Handlers
func (p *Pipeline) Start() {
	var input chan interface{}
	var output chan interface{}

	for index := range p.handlers {
		if index == 0 {
			input = p.Input()
		} else {
			input = output
		}

		if index == len(p.handlers)-1 {
			output = p.Output()
		} else {
			output = p.NewChannel()
		}

		p.startHandlerByIndex(index, input, output)
	}
}

// Wait holds thread until all handlers stops
func (p *Pipeline) Wait() {
	p.wg.Wait()
}

// Stop stops all gourotines and close channels
func (p *Pipeline) Stop() {
	close(p.Input())
}

func (p *Pipeline) startHandlerByIndex(index int, input, output chan interface{}) {
	scale := p.scale[index]
	p.wg.Add(scale)

	for i := 0; i < scale; i++ {
		go p.handlers[index].Handle(input, output, &p.wg)
	}
}
