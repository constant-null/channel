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
func (p *Pipeline) Add(handler Handler) {
	p.handlers = append(p.handlers, handler)
}

// Start starting Handlers
func (p *Pipeline) Start() {
	input := p.NewChannel()
	for index, handler := range p.handlers {
		p.wg.Add(1)
		switch index {
		case 0:
			go handler.Handle(p.Input(), input, &p.wg)
		case len(p.handlers) - 1:
			go handler.Handle(input, p.Output(), &p.wg)
		default:
			output := p.NewChannel()
			go handler.Handle(input, output, &p.wg)
			input = output
		}
	}
}

// Wait holds thread until all handlers stops
func (p *Pipeline) Wait() {
	p.wg.Wait()
}