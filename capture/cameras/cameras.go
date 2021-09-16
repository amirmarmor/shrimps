package cameras

import (
	"encoding/json"
	"fmt"
	"gocv.io/x/gocv"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"time"
	"www.seawise.com/shrimps/common/core"
	"www.seawise.com/shrimps/common/exposed"
	"www.seawise.com/shrimps/common/log"
)

type Channel struct {
	numCameras int
	offset     int
	counter    int
	config     *core.Configuration
	Captures   []*gocv.VideoCapture
	Images     []gocv.Mat
	Writers    []*gocv.VideoWriter
	Windows    []*gocv.Window
	Recording	 map[string]time.Time
	ticker     *time.Ticker
}

func inArray(needle int, haystack []int) bool {
	for _, val := range haystack {
		if needle == val {
			return true
		}
	}
	return false
}

func Create() (*Channel, error) {
	channel := &Channel{}
	err := channel.GetConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get config: %v", err)
	}
	return channel, nil
}

func (c *Channel) closeCameras() {
	for i := 0; i < len(c.Captures); i++ {
		c.Captures[i].Close()
		c.Images[i].Close()
		c.Writers[i].Close()
		c.Windows[i].Close()
	}
}

func (c *Channel) Init(saveFileName string) error {
	c.closeCameras()
	for camera := c.offset; camera < (c.numCameras * 2); camera += 2 {
		vc, err := gocv.OpenVideoCapture(camera)
		if err != nil {
			return fmt.Errorf("failed to capture video %v: ", err)
		}
		c.Captures = append(c.Captures, vc)

		img := gocv.NewMat()
		c.Images = append(c.Images, img)

		window := gocv.NewWindow(fmt.Sprintf("channel-%v", camera))
		c.Windows = append(c.Windows, window)
		ok := vc.Read(&img)
		if !ok {
			return fmt.Errorf("failed to read")
		}

		saveFileName = saveFileName + "-" + strconv.Itoa(camera) + ".avi"

		writer, err := gocv.VideoWriterFile(saveFileName, "MJPG", 25, img.Cols(), img.Rows(), true)
		if err != nil {
			return fmt.Errorf("failed to create writer", err)
		}
		c.Writers = append(c.Writers, writer)
	}
	return nil
}

func (c *Channel) Start() error {
	c.ticker = time.NewTicker(time.Second * 1)
	for {
		select {
		case <-c.ticker.C:
			c.update()
		default:
			c.capture()
		}
	}
}

func (c* Channel) update() error{
	fmt.Println("go ticker")
	err := c.GetConfig()
	if err != nil {
		return err
	}
	c.capture()
	return nil
}

func (c *Channel) capture() error {
	for i := 0; i < len(c.Captures); i++ {
		ok := c.Captures[i].Read(&c.Images[i])
		if !ok {
			return fmt.Errorf("channel closed %v\n", i)
		}

		if c.Images[i].Empty() {
			continue
		}

		write, err := c.check(i)
		if err != nil {
			return err
		}
		if write {
			c.Writers[i].Write(c.Images[i])
		}

		if inArray(i, c.config.Show) {
			c.Windows[i].IMShow(c.Images[i])
			c.Windows[i].WaitKey(1)
		}
	}
	return nil
}

func (c *Channel) GetConfig() error {
	now := time.Now()
	c.counter++
	response, err := http.Get("http://127.0.0.1:1323/config")
	if err != nil {
		return fmt.Errorf("failed to get configuration: %v", err)
	}

	configJson, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("failed to read body: %v", err)
	}

	config := &core.Configuration{}
	err = json.Unmarshal(configJson, config)
	if err != nil {
		return fmt.Errorf("failed to unmarshal config: %v", err)
	}

	numCameras, err := strconv.Atoi(config.Cameras)
	if err != nil {
		return fmt.Errorf("failed to convert configuration: %v", err)
	}

	offset, err := strconv.Atoi(config.Offset)
	if err != nil {
		return fmt.Errorf("failed to convert configuration: %v", err)
	}

	path, err := createSavePath()
	if err != nil {
		return fmt.Errorf("failed to create path: %v", err)
	}

	saveFileName := path + "/" +
		strconv.Itoa(now.Hour()) +
		strconv.Itoa(now.Minute()) +
		strconv.Itoa(now.Second())

	if len(c.Captures) == 0 ||
		c.numCameras != numCameras ||
		c.offset != offset {
		c.numCameras = numCameras
		c.offset = offset
		c.config = config
		c.Init(saveFileName)
	}
	log.V5(strconv.Itoa(c.counter))

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

func (c *Channel) check(camera int) (bool, error) {
	if inArray(camera, c.config.Record){
		return true, nil
	}
	return c.checkRule()

}

func (c *Channel) checkRule() (bool, error) {
	now := time.Now()
	zeroTime, err := time.Parse(time.RFC3339, exposed.ZeroTime)
	if err != nil {
		return false, err
	}

	for _, rule := range c.config.Rules {
		start, err  := strconv.Atoi(rule.Start)
		if err != nil {
			return false, fmt.Errorf("failed to convert rule: %v", err)
		}
		duration, err := strconv.Atoi(rule.Duration)
		if err != nil {
			return false, fmt.Errorf("failed to convert rule: %v", err)
		}

		bar := GetTimeField(rule.Recurring)

		if start == bar {
			if c.Recording[rule.Id] == zeroTime {
				c.Recording[rule.Id] = now
				return true, nil
			}

			if now.Sub(c.Recording[rule.Id]) <= time.Second*time.Duration(duration) {
				return true, nil
			}
		} else {
			c.Recording[rule.Id] = zeroTime
		}
	}
	return false, nil
}

func GetTimeField(s string) int {
	now := time.Now()
	r := reflect.ValueOf(now)
	f := reflect.Indirect(r).FieldByName(s)
	return int(f.Int())
}

