package capture

import (
	"fmt"
	"time"
	"www.seawise.com/shrimps/backend/core"
	"www.seawise.com/shrimps/backend/log"
)

type Capture struct {
	counter     int
	run         bool
	manager     *core.ConfigManager
	Channels    []*Channel
	Recording   map[string]time.Time
	Action      chan *ShowRecord
	StopChannel chan string
	Errors      chan error
}

type ShowRecord struct {
	Type    string
	Channel int
}

func Create(config *core.ConfigManager) *Capture {
	return &Capture{
		manager:     config,
		run:         true,
		Action:      make(chan *ShowRecord, 0),
		StopChannel: make(chan string, 0),
	}
}

func (c *Capture) Init() error {
	return c.detectCameras()
}

func (c *Capture) detectCameras() error {
	isDetecting := true
	i := c.manager.Config.Offset
	c.Channels = make([]*Channel, 0)
	for isDetecting {
		channel := CreateChannel(i, c.manager.Config.Rules)
		err := channel.Init()
		if err != nil {
			if i > 2 {
				isDetecting = false
			} else {
				i+=2
				continue
			}
		} else {
			c.Channels = append(c.Channels, channel)
			i+=2
		}
	}
	return nil
}

func (c *Capture) Start() {
	for c.run {
		select {
		case action := <-c.Action:
			c.update(action)
		case s := <-c.StopChannel:
			c.stop(s)
		default:
			c.capture()
		}
	}
	c.StopChannel <- "restarting"
}

func (c *Capture) Update() {
	c.StopChannel <- "stopping"
	for !c.run {
		select {
		case s := <-c.StopChannel:
			c.restart(s)
		default:
			continue
		}
	}
}

func (c *Capture) stop(s string) {
	log.V5(fmt.Sprintf("capture - %s", s))
	c.run = false
}

func (c *Capture) restart(s string) error {
	log.V5(fmt.Sprintf("capture - %s", s))
	c.run = true

	for i := 0; i < len(c.Channels); i++ {
		c.Channels[i].close()
	}

	c.Channels = make([]*Channel, 0)

	err := c.Init()
	if err != nil {
		return err
	}

	go c.Start()
	return nil
}

func (c *Capture) update(action *ShowRecord) error {
	if action.Type == "show" {
		c.Channels[action.Channel].Show = !c.Channels[action.Channel].Show
	}
	if action.Type == "record" {
		c.Channels[action.Channel].Record = !c.Channels[action.Channel].Record
	}

	if action.Type == "config" {
		err := c.manager.GetConfig()
		for _, channel := range c.Channels {
			channel.rules = c.manager.Config.Rules
		}
		if err != nil {
			return err
		}
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
