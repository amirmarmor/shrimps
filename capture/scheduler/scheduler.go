package scheduler

import (
	"fmt"
	"strconv"
	"time"
	"www.seawise.com/shrimps/capture/channels"
	"www.seawise.com/shrimps/common/core"
)

type Scheduler struct {
	Manager  *core.ConfigManager
	Recorder *channels.Recorder
	ticker   *time.Ticker
	errors   chan error
}

func Create(manager *core.ConfigManager, recorder *channels.Recorder) *Scheduler {
	scheduler := &Scheduler{
		Manager:  manager,
		Recorder: recorder,
		ticker:   time.NewTicker(time.Second * 1),
		errors:   make(chan error, 0),
	}

	return scheduler
}

func (s *Scheduler) Start() error {
	for {
		select {
		case <-s.ticker.C:
			s.doWork()
		}
	}
}

func (s *Scheduler) doWork() error {
	now := time.Now()
	config, err := s.Manager.GetConfig()
	if err != nil {
		return fmt.Errorf("failed to get configuration: %v", err)
	}

	cameras, err := strconv.Atoi(config.Cameras)
	if err != nil {
		return fmt.Errorf("failed to convert configuration: %v", err)
	}

	offset, err := strconv.Atoi(config.Offset)
	if err != nil {
		return fmt.Errorf("failed to convert configuration: %v", err)
	}

	for _, rule := range config.Rules {
		if rule.Recurring == "h" {
			err := s.checkRule(now.Second(), now, rule, cameras, offset)
			if err != nil {
				fmt.Println(err)
			}
		}
		if rule.Recurring == "d" {
			err := s.checkRule(now.Hour(), now, rule, cameras, offset)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
	return nil
}

func (s *Scheduler) checkRule(bar int, now time.Time, rule core.Rule, cameras int, offset int) error {
	start, err := strconv.Atoi(rule.Start)
	if err != nil {
		return err
	}
	fmt.Println(start, bar, rule.Duration)
	if start == bar {
		duration, err := strconv.Atoi(rule.Duration)
		if err != nil {
			return err
		}
		s.Recorder.Start = now
		s.Recorder.Duration = time.Second * time.Duration(duration)
		s.Recorder.Cameras = cameras
		s.Recorder.Offset = offset

		err = s.Recorder.Record()
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	return nil
}
