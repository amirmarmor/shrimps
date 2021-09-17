package action

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"www.seawise.com/shrimps/backend/capture"
)

type Executor struct {
	Capt *capture.Capture
}

type ActionRequest struct {
	Type    string `json:"type"`
	Channel string `json:"channel"`
}
type SetResponse struct {
	Message string
}

func Produce(capt *capture.Capture) *Executor {
	executor := Executor{
		Capt: capt,
	}
	return &executor
}

func (executor *Executor) Do(c echo.Context) error {
	actionRequest := &ActionRequest{}
	err := c.Bind(actionRequest)
	if err != nil {
		return fmt.Errorf("failed to bind: %v", err)
	}

	channel, err := strconv.Atoi(actionRequest.Channel)
	if err != nil {
		return fmt.Errorf("failed to convert channel: %v", err)
	}

	action := &capture.ShowRecord{

		Type: actionRequest.Type,
		Channel: channel,
	}

	executor.Capt.Action <- action

	return c.JSON(http.StatusOK, &SetResponse{Message: "done"})
}
