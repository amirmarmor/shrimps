package entrypoint

import (
	"www.seawise.com/shrimps/common/core"
	"www.seawise.com/shrimps/common/exposed"
	"www.seawise.com/shrimps/common/persistance"
	"www.seawise.com/shrimps/guibackend/executors/configuration"
	"www.seawise.com/shrimps/guibackend/web"
)

type EntryPoint struct {
	persistance *persistance.Persist
	web         *web.Web
	manager     *core.ConfigManager
}

func (p *EntryPoint) Run() {
	p.buildBlocks()
	p.addExecutors()

	p.web.Start()
}

func (p *EntryPoint) buildBlocks() {
	persist, err := persistance.Create()
	if err != nil {
		panic(err)
	}

	p.persistance = persist
	p.web = web.Create()

	p.manager, err = core.Produce(persist)
	if err != nil {
		panic(err)
	}
}

func (p *EntryPoint) addExecutors() {

	configurationExecutor := configuration.Produce(p.manager)
	p.web.Client.GET(exposed.GetConfigUrl, configurationExecutor.GetConfig)
	p.web.Client.POST(exposed.SetConfigUrl, configurationExecutor.SetConfig)
}
