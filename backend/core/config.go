package core

import (
	"encoding/json"
	"fmt"
	"www.seawise.com/shrimps/backend/exposed"
	"www.seawise.com/shrimps/backend/persistance"
)

type Configuration struct {
	Cameras string   `json:"cameras"`
	Offset  string   `json:"offset"`
	Rules   []Rule   `json:"rules"`
}

type ConfigManager struct {
	Persist *persistance.Persist
	Config  *Configuration
}

type Rule struct {
	Id        string `json:"id"`
	Recurring string `json:"recurring"`
	Start     string `json:"start"`
	Duration  string `json:"duration"`
}

func Produce(persistanceApi *persistance.Persist) (*ConfigManager, error) {
	InitFlags()

	manager := ConfigManager{
		Persist: persistanceApi,
	}

	config, err := manager.GetConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get configuration: %v", err)
	}

	manager.Config = config

	return &manager, nil
}

func (cm *ConfigManager) GetConfig() (*Configuration, error) {
	configJson, err := cm.Persist.Get(exposed.ConfigKey)
	if err != nil {
		return nil, err
	}

	config := &Configuration{}
	if configJson != "" {
		err = json.Unmarshal([]byte(configJson), &config)
		if err != nil {
			return nil, err
		}
		return config, nil
	}

	config.Cameras = Defaults.Cameras
	config.Offset = Defaults.Offset

	rules := make([]Rule, 0)
	err = json.Unmarshal([]byte(Defaults.Rules), &rules)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %v", err)
	}

	show := make([]int, 0)
	err = json.Unmarshal([]byte(Defaults.Show), &show)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %v", err)
	}

	record := make([]int, 0)
	err = json.Unmarshal([]byte(Defaults.Record), &record)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %v", err)
	}

	config.Rules = rules

	return config, nil
}

func (cm *ConfigManager) SetConfig(configJson string) error {
	err := cm.Persist.Set(exposed.ConfigKey, configJson)
	if err != nil {
		return fmt.Errorf("failed to set: %v", err)
	}

	config := &Configuration{}
	err = json.Unmarshal([]byte(configJson), config)
	if err != nil {
		return fmt.Errorf("failed to unmarshal: %v", err)
	}
	cm.Config = config
	return nil
}

//func (cm *ConfigManager) GetRules() (*Rules, error) {
//	rules := &Rules{}
//	err := json.Unmarshal([]byte(cm.Config.Rules), rules)
//	if err != nil {
//		return nil, fmt.Errorf("failed to unmarshal rules: %v", err)
//	}
//	return rules, nil
//}
