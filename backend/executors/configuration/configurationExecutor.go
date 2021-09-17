package configuration

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"www.seawise.com/shrimps/backend/capture"
	"www.seawise.com/shrimps/backend/core"
)

type Executor struct {
	Manager *core.ConfigManager
	Capt    *capture.Capture
}

type SetResponse struct {
	Message string
}

type GetResponse struct {
	Config *core.Configuration `json:"config"`
	Show   []int               `json:"show"`
	Record []int               `json:"record"`
}

func Produce(configManager *core.ConfigManager, capt *capture.Capture) *Executor {
	executor := Executor{
		Manager: configManager,
		Capt:    capt,
	}
	return &executor
}

func (executor *Executor) GetConfig(c echo.Context) error {
	show := make([]int, 0)
	record := make([]int, 0)
	for i, channel := range executor.Capt.Channels {
		if channel.Show {
			show = append(show, i)
		}
		if channel.Record {
			record = append(record, i)
		}
	}
	response := &GetResponse{
		executor.Manager.Config,
		show,
		record,
	}
	return c.JSON(http.StatusOK, response)
}

func (executor *Executor) SetConfig(c echo.Context) error {
	config := &core.Configuration{}

	err := c.Bind(config)
	if err != nil {
		return fmt.Errorf("failed to bind: %v", err)
	}

	configJson, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal: %v", err)
	}

	temp := string(configJson)
	err = executor.Manager.SetConfig(temp)
	if err != nil {
		return fmt.Errorf("failed to set configuration: %v", err)
	}

	return c.JSON(http.StatusOK, &SetResponse{Message: "done"})
}
