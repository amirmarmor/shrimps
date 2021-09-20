package capture

import (
	"fmt"
	"time"
	"www.seawise.com/shrimps/backend/core"
)

type Capture struct {
	counter    int
	config     *core.Configuration
	Channels   []*Channel
	Recording  map[string]time.Time
	Action     chan *ShowRecord
	Errors     chan error
}

type ShowRecord struct {
	Type    string
	Channel int
}

func Create(config *core.Configuration) *Capture {
	return &Capture{
		config: config,
		Action: make(chan *ShowRecord, 0),
	}
}

func (c *Capture) Init() error {
	return c.detectCameras()
}

func (c *Capture) detectCameras() error {
	detecting := true
	i := c.config.Offset
	for detecting {
		channel := CreateChannel(i, c.config.Rules)
		err := channel.Init()
		if err != nil {
			detecting = false
		} else {
			c.Channels = append(c.Channels, channel)
			i += 2
		}
	}

	return nil
}

func (c *Capture) Start() {
	for true {
		select {
		case action := <-c.Action:
			c.update(action)
		default:
			c.capture()
		}
	}
}

//func (c *Capture) Update(config *core.Configuration) error {
//	for i := 0; i < len(c.Channels); i++ {
//		c.Channels[i].close()
//	}
//	c.config = config
//	err := c.Init()
//	if err != nil {
//		return err
//	}
//	go c.Start()
//	return nil
//}

func (c *Capture) update(action *ShowRecord) error {
	if action.Type == "show" {
		c.Channels[action.Channel].Show = !c.Channels[action.Channel].Show
	}
	if action.Type == "record" {
		c.Channels[action.Channel].Record = !c.Channels[action.Channel].Record
	}
	err := c.capture()
	if err != nil {
		return err
	}
	return nil
}

func (c *Capture) capture() error {
	for _, channel := range c.Channels {
		err := channel.Read()
		if err != nil {
			return fmt.Errorf("capture failed: %v", err)
		}
	}
	return nil
}

