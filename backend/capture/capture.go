package capture

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
	"www.seawise.com/shrimps/backend/core"
	"www.seawise.com/shrimps/backend/exposed"
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
	for detecting {
		i := c.config.Offset
		channel := CreateChannel(i)
		err := channel.Init()
		if err != nil {
			detecting = false
		} else {
			c.Channels = append(c.Channels, channel)
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

func (c *Capture) checkRule() bool {
	now := time.Now()
	zeroTime, err := time.Parse(time.RFC3339, exposed.ZeroTime)
	if err != nil {
		c.Errors <- err
		return false
	}

	for _, rule := range c.config.Rules {
		start, err := strconv.Atoi(rule.Start)
		if err != nil {
			c.Errors <- fmt.Errorf("failed to convert rule: %v", err)
			return false
		}
		duration, err := strconv.Atoi(rule.Duration)
		if err != nil {
			c.Errors <- fmt.Errorf("failed to convert rule: %v", err)
			return false
		}

		bar := GetTimeField(rule.Recurring)

		if start == bar {
			if c.Recording[rule.Id] == zeroTime {
				c.Recording[rule.Id] = now
				return true
			}

			if now.Sub(c.Recording[rule.Id]) <= time.Second*time.Duration(duration) {
				return true
			}
		} else {
			c.Recording[rule.Id] = zeroTime
		}
	}
	return false
}

func GetTimeField(s string) int {
	now := time.Now()
	r := reflect.ValueOf(now)
	f := reflect.Indirect(r).FieldByName(s)
	return int(f.Int())
}
