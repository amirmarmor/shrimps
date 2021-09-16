package persistance

import (
	"encoding/json"
	"github.com/namsral/flag"
	"www.seawise.com/shrimps/common/log"
)

type RedisConfig struct {
	Host string
}

var Config RedisConfig

func InitFlags() {
	flag.StringVar(&Config.Host, "host", "redis", "redis host")

	log.AddNotify(postParse)
}

func postParse() {
	marshal, err := json.Marshal(Config)
	if err != nil {
		log.Fatal("marshal config failed: %v", err)
	}

	log.V5("configuration loaded: %v", string(marshal))
}
