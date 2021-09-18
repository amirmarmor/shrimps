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

func Produce(channel int) (*Channel, error) {
	now := time.Now()
	vc, err := gocv.OpenVideoCapture(channel)
	if err != nil {
		return nil, fmt.Errorf("failed to capture video %v: ", err)
	}

	img := gocv.NewMat()

	ok := vc.Read(&img)
	if !ok {
		return nil, fmt.Errorf("failed to read")
	}

	path, err := createSavePath()
	if err != nil {
		return nil, fmt.Errorf("failed to create path: %v", err)
	}

	stream := mjpeg.NewStream()

	window := gocv.NewWindow("channel-" + strconv.Itoa(channel))
	window.ResizeWindow(1,1)

	saveFileName := path + "/" +
		strconv.Itoa(now.Hour()) +
		strconv.Itoa(now.Minute()) +
		strconv.Itoa(now.Second()) +
		"-" + strconv.Itoa(channel) +
		".avi"

	writer, err := gocv.VideoWriterFile(saveFileName, "MJPG", 25, img.Cols(), img.Rows(), true)
	if err != nil {
		return nil, fmt.Errorf("failed to create writer", err)
	}

	c := &Channel{}
	c.name = channel
	c.cap = vc
	c.image = img
	c.writer = writer
	c.init = true
	c.Stream = stream
	c.Window = window

	return c, nil

}

func (c *Channel) close() {
	c.cap.Close()
	c.image.Close()
	c.writer.Close()
	c.Window.Close()
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
