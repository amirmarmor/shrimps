package channels

import (
	"fmt"
	"gocv.io/x/gocv"
	"os"
	"strconv"
	"time"
	"www.seawise.com/shrimps/common/log"
)

type Recorder struct {
	Start    time.Time
	Duration time.Duration
	path     string
	Cameras  int
	Offset   int
	Recording bool
}

func Create() (*Recorder, error) {
	path, err := createSavePath()
	if err != nil {
		return nil, err
	}

	recorder := &Recorder{
		time.Now(),
		0,
		path,
		0,
		0,
		false,
	}

	return recorder, nil
}

func (r *Recorder) Record() error {
	if !r.Recording {
		now := time.Now()
		cameras := make([]*gocv.VideoCapture, 0)
		images := make([]gocv.Mat, 0)
		writers := make([]*gocv.VideoWriter, 0)

		for camera := r.Offset; camera < (r.Cameras * 2); camera += 2 {
			vc, err := gocv.OpenVideoCapture(camera)
			if err != nil {
				return fmt.Errorf("failed to capture video %v: ", err)
			}
			cameras = append(cameras, vc)

			img := gocv.NewMat()
			images = append(images, img)

			ok := vc.Read(&img)
			if !ok {
				return fmt.Errorf("failed to read")
			}

			saveFileName := r.path + "/" +
				strconv.Itoa(camera) + "-" +
				strconv.Itoa(now.Hour()) + "-" +
				strconv.Itoa(now.Minute()) + "-" +
				strconv.Itoa(now.Second()) + ".avi"

			writer, err := gocv.VideoWriterFile(saveFileName, "MJPG", 25, img.Cols(), img.Rows(), true)
			if err != nil {
				return fmt.Errorf("failed to create writer", err)
			}
			writers = append(writers, writer)
		}

		for now.Sub(r.Start) <= r.Duration {
			now = time.Now()
			fmt.Println("1111", now, now.Sub(r.Start), r.Duration)
			for i := 0; i < len(cameras); i++ {
				if ok := cameras[i].Read(&images[i]); !ok {
					return fmt.Errorf("channel closed %v\n", i)
				}

				if images[i].Empty() {
					continue
				}

				writers[i].Write(images[i])
			}
		}
		fmt.Println("aaaaaaa")
		r.Duration = 0
		r.Start = now
		for i := 0; i < len(cameras); i++ {
			cameras[i].Close()
			images[i].Close()
			writers[i].Close()
		}
		return nil
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

