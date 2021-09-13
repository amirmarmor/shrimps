package entrypoint

import (
	"www.seawise.com/shrimps/capture/channels"
	"www.seawise.com/shrimps/capture/scheduler"
	"www.seawise.com/shrimps/common/core"
	"www.seawise.com/shrimps/common/persistance"
)

type EntryPoint struct {
	persistance *persistance.Persist
	manager     *core.ConfigManager
	recorder    *channels.Recorder
	scheduler 	*scheduler.Scheduler
}

func (p *EntryPoint) Run() {
	p.buildBlocks()

	p.scheduler.Start()
}

func (p *EntryPoint) buildBlocks() {
	persist, err := persistance.Create()
	if err != nil {
		panic(err)
	}

	p.persistance = persist

	p.manager, err = core.Produce(persist)
	if err != nil {
		panic(err)
	}

	p.recorder, err = channels.Create()
	if err != nil {
		panic(err)
	}

	p.scheduler = scheduler.Create(p.manager, p.recorder)

}

