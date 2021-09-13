package configuration

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"www.seawise.com/shrimps/common/core"
)

type Executor struct {
	Manager *core.ConfigManager
}

type SetResponse struct {
	Message string
}

func Produce(configManager *core.ConfigManager) *Executor {
	executor := Executor{
		Manager: configManager,
	}
	return &executor
}

func (executor *Executor) GetConfig(c echo.Context) error {
	return c.JSON(http.StatusOK, executor.Manager.Config)
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

	err = executor.Manager.SetConfig(string(configJson))
	if err != nil {
		return fmt.Errorf("failed to set configuration: %v", err)
	}

	return c.JSON(http.StatusOK, &SetResponse{Message: "done"})
}
