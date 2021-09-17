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
	numCameras, err := strconv.Atoi(c.config.Cameras)
	if err != nil {
		return fmt.Errorf("failed to convert config: %v", err)
	}

	offset, err := strconv.Atoi(c.config.Offset)
	if err != nil {
		return fmt.Errorf("failed to convert config: %v", err)
	}

	for i := offset; i < numCameras * 2; i += 2 {
		channel, err := Produce(i)
		if err != nil {
			return fmt.Errorf("failed to produce channel: %v", err)
		}
		c.Channels = append(c.Channels, channel)
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

func (c *Capture) Update(config *core.Configuration) error {
	for i := 0; i < len(c.Channels); i++ {
		c.Channels[i].close()
	}
	c.config = config
	err := c.Init()
	if err != nil {
		return err
	}
	go c.Start()
	return nil
}

func (c *Capture) update(action *ShowRecord) {
	if action.Type == "show" {
		c.Channels[action.Channel].Show = !c.Channels[action.Channel].Show
		if c.Channels[action.Channel].Show {
			c.Channels[action.Channel].createWindow()
		} else {
			c.Channels[action.Channel].window.Close()
		}
	}
	if action.Type == "record" {
		c.Channels[action.Channel].Record = !c.Channels[action.Channel].Record
	}
	c.capture()
}

func (c *Capture) capture() {
	//log.V5("capture")
	c.counter ++
	//if c.counter > 1000000000 {
	//	panic("000000000")
	//}
	for i := 0; i < len(c.Channels); i++ {
		if !c.Channels[i].Show && !c.Channels[i].Record {
			continue
		}

		if !c.Channels[i].init {
			c.Start()
		}

		ok := c.Channels[i].cap.Read(&c.Channels[i].image)
		if !ok {
			c.Errors <- fmt.Errorf("channel closed %v\n", i)
		}

		if c.Channels[i].image.Empty() {
			continue
		}

		if c.Channels[i].Record {
			err := c.Channels[i].writer.Write(c.Channels[i].image)
			if err != nil {
				c.Errors <- fmt.Errorf("failed to Write: %v", err)
			}
		}

		if c.Channels[i].Show {
			c.Channels[i].window.IMShow(c.Channels[i].image)
			c.Channels[i].window.WaitKey(1)
		}
	}
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
