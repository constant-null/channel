package balancer

// EmptyChannelListError occures when there is no channels has been added to balancer
type EmptyChannelListError struct {
}

func (err EmptyChannelListError) Error() string {
	return "At least one channel should be added"
}

// Balancer helps balance incoming tagged data between fixed amount of channels.
// Tags are unique across all channels, and one channel can have several tags.
//
// It fully utilizes channels only if there are more tags than channels
type Balancer struct {
	channels    [](chan interface{})
	tags        map[string]int
	nextChannel int
}

// New creates new channel balancer
func New() *Balancer {
	return &Balancer{tags: make(map[string]int)}
}

// Add adds new channel to registy
func (c *Balancer) Add(channel chan interface{}) {
	c.channels = append(c.channels, channel)
}

// Get returns next channel from balancer
func (c *Balancer) Get(tag string) (chan interface{}, error) {
	if len(c.channels) == 0 {
		return nil, EmptyChannelListError{}
	}

	channelID, ok := c.tags[tag]

	if !ok {
		channelID = c.getNextChannel()
		c.tags[tag] = channelID
	}

	return c.channels[channelID], nil
}

// Close closes all channels
func (c *Balancer) Close() {
	for _, channel := range c.channels {
		close(channel)
	}
	c.channels = nil
	c.tags = make(map[string]int)
}

func (c *Balancer) getNextChannel() int {
	id := c.nextChannel

	if c.nextChannel+1 == len(c.channels) {
		c.nextChannel = 0
	} else {
		c.nextChannel++
	}

	return id
}
