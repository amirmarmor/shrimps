package entrypoint

import (
	"www.seawise.com/shrimps/capture/cameras"
)

type EntryPoint struct {
	recorder *cameras.Channel
	channel  *cameras.Channel
}

func (p *EntryPoint) Run() {
	p.buildBlocks()

	p.channel.Start()
}

func (p *EntryPoint) buildBlocks() {
	var err error
	p.channel, err = cameras.Create()
	if err != nil {
		panic(err)
	}
}
