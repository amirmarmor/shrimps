package core

import (
	"encoding/json"
	"github.com/namsral/flag"
	"www.seawise.com/shrimps/common/log"
)

type Default = struct {
	Cameras string
	Offset  string
	Rules   string
}

var Defaults Default

func InitFlags() {
	defaultRule := `[{"recurring": "h", "start": "99", "duration": "180"}, {"recurring": "d", "start": "15", "duration": "60"}]`

	flag.StringVar(&Defaults.Cameras, "cameras", "1", "The number of cameras connected")
	flag.StringVar(&Defaults.Offset, "offset", "0", "The offset of the first camera - if no webcam = 0")
	flag.StringVar(&Defaults.Rules, "rules", defaultRule, "rules")

	log.AddNotify(postParse)
}

func postParse() {
	marshal, err := json.Marshal(Defaults)
	if err != nil {
		log.Fatal("marshal config failed: %v", err)
	}

	log.V5("configuration loaded: %v", string(marshal))
}
