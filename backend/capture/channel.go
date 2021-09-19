package capture

import (
	"fmt"
	"gocv.io/x/gocv"
	"os"
	"strconv"
	"time"
	"www.seawise.com/shrimps/backend/log"
	"www.seawise.com/shrimps/backend/mjpeg"
)

type Channel struct {
	name   int
	init   bool
	cap    *gocv.VideoCapture
	image  gocv.Mat
	writer *gocv.VideoWriter
	Stream *mjpeg.Stream
	Window *gocv.Window
	Show   bool
	Record bool
}

func CreateChannel(channel int) *Channel {
	return &Channel{
		name: channel,
		Stream: mjpeg.NewStream(),
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

	path, err := createSavePath()
	if err != nil {
		return fmt.Errorf("failed to create path: %v", err)
	}

	window := gocv.NewWindow("channel-" + strconv.Itoa(c.name))
	window.ResizeWindow(1,1)

	saveFileName := path + "/" +
		strconv.Itoa(now.Hour()) +
		strconv.Itoa(now.Minute()) +
		strconv.Itoa(now.Second()) +
		"-" + strconv.Itoa(c.name) +
		".avi"

	writer, err := gocv.VideoWriterFile(saveFileName, "MJPG", 25, img.Cols(), img.Rows(), true)
	if err != nil {
		return fmt.Errorf("failed to create writer", err)
	}

	c.cap = vc
	c.image = img
	c.writer = writer
	c.init = true
	c.Window = window

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

	err =c.writer.Close()
	if err != nil {
		return fmt.Errorf("failed to close writer: %v", err)
	}

	err = c.Window.Close()
	if err != nil {
		return fmt.Errorf("failed to close window: %v", err)
	}

	c.init = false
	return nil
}

func (c *Channel) Read() error {
	if !c.Show && !c.Record {
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

	if c.Record {
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
	if c.Show {
		c.Window.IMShow(c.image)
		c.Window.WaitKey(1)
	}
	 return nil
}

func createSavePath() (string, error) {
	_, err := os.Stat("videos")

	if os.IsNotExist(err) {
		log.V5("videos directory doesnt exist. creating it now!")
		err := os.Mkdir("videos", 0777)
		if err != nil {
			log.Error("couldnt create images directory", err)
			return "", err
		}
	}

	now := time.Now()
	y, m, d := now.Date()
	path := fmt.Sprintf("videos/%v", fmt.Sprintf("%v-%v-%v", y, m, d))
	_, err = os.Stat(path)

	if !os.IsNotExist(err) {
		err := os.RemoveAll(path)
		if err != nil {
			log.Error("couldnt remove folder", path)
		}
	}

	log.V5("creating file direcotry!")
	err = os.Mkdir(path, 0777)
	if err != nil {
		log.Error("couldnt create images directory", err)
		return "", err
	}
	return path, nil
}
