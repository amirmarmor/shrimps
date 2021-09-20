package capture

import (
	"fmt"
	"gocv.io/x/gocv"
	"os"
	"reflect"
	"time"
	"www.seawise.com/shrimps/backend/core"
	"www.seawise.com/shrimps/backend/log"
	"www.seawise.com/shrimps/backend/mjpeg"
)

type Channel struct {
	name      int
	init      bool
	cap       *gocv.VideoCapture
	image     gocv.Mat
	writer    *gocv.VideoWriter
	Stream    *mjpeg.Stream
	Show      bool
	Record    bool
	Recordings map[int64]*Recording
	rules     []core.Rule
}

type Recording struct {
	isRecording bool
	startTime   time.Time
}

func CreateChannel(channel int, rules []core.Rule) *Channel {
	return &Channel{
		name:   channel,
		Stream: mjpeg.NewStream(),
		rules:  rules,
		Recordings: make(map[int64]*Recording),
	}
}

func (c *Channel) Init() error {
	now := time.Now()
	vc, err := gocv.OpenVideoCapture(c.name)
	if err != nil {
		return fmt.Errorf("Init failed to capture video %v: ", err)
	}

	img := gocv.NewMat()

	ok := vc.Read(&img)
	if !ok {
		return fmt.Errorf("Init failed to read")
	}

	path, err := c.createSavePath()
	if err != nil {
		return fmt.Errorf("failed to create path: %v", err)
	}

	saveFileName := path + "/" + now.Format("2006-01-02--15-04-05") + ".avi"

	writer, err := gocv.VideoWriterFile(saveFileName, "MJPG", 25, img.Cols(), img.Rows(), true)
	if err != nil {
		return fmt.Errorf("failed to create writer", err)
	}

	c.cap = vc
	c.image = img
	c.writer = writer
	c.init = true

	return nil
}

func (c *Channel) close() error {
	err := c.cap.Close()
	if err != nil {
		return fmt.Errorf("failed to close capture: %v", err)
	}
	err = c.image.Close()
	if err != nil {
		return fmt.Errorf("failed to close image: %v", err)
	}

	err = c.writer.Close()
	if err != nil {
		return fmt.Errorf("failed to close writer: %v", err)
	}

	c.init = false
	return nil
}

func (c *Channel) Read() error {
	if !c.Show && !c.Record && !c.checkRules(){
		if c.init {
			err := c.close()
			if err != nil {
				return fmt.Errorf("read failed to close: %v", err)
			}
		}
		return nil
	}

	if !c.init {
		err := c.Init()
		if err != nil {
			return fmt.Errorf("read init failed to close: %v", err)
		}

	}

	ok := c.cap.Read(&c.image)
	if !ok {
		return fmt.Errorf("read encountered channel closed %v\n", c.name)
	}

	if c.image.Empty() {
		return nil
	}

	if c.Record || c.checkRules() {
		err := c.writer.Write(c.image)
		if err != nil {
			return fmt.Errorf("read failed to write: %v", err)
		}
	}

	buffer, err := gocv.IMEncode(".jpg", c.image)
	if err != nil {
		return fmt.Errorf("read failed to encode: %v", err)
	}

	c.Stream.UpdateJPEG(buffer.GetBytes())
	return nil
}

func (c *Channel) createSavePath() (string, error) {
	_, err := os.Stat("videos")

	if os.IsNotExist(err) {
		log.V5("videos directory doesnt exist. creating it now!")
		err := os.Mkdir("videos", 0777)
		if err != nil {
			log.Error("couldnt create images directory", err)
			return "", err
		}
	}

	path := fmt.Sprintf("videos/channel-%v", c.name)
	_, err = os.Stat(path)

	if os.IsNotExist(err) {
		log.V5("creating file direcotry!")
		err = os.Mkdir(path, 0777)
		if err != nil {
			log.Error("couldnt create images directory", err)
			return "", err
		}
	}

	return path, nil
}

func (c *Channel) checkRules() bool {
	now := time.Now()
	for _, rule := range c.rules {
		bar := GetTimeField(rule.Recurring, now)
		if rule.Start == bar {
			if c.Recordings[rule.Id] == nil {
				c.Recordings[rule.Id] = &Recording{
					true,
					now,
				}
				return true
			}
		}

		if c.Recordings[rule.Id] != nil && now.Sub(c.Recordings[rule.Id].startTime) <= time.Second*time.Duration(rule.Duration) {
			return true
		}

		c.Recordings[rule.Id] = nil
	}
	return false
}

func GetTimeField(s string, now time.Time) int64 {
	r := reflect.ValueOf(now).MethodByName(s)
	f := r.Call(nil)
	return int64(f[0].Interface().(int))
}
