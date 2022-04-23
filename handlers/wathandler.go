package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/USACE/filestore"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/labstack/echo/v4"
	"github.com/usace/wat-api/utils"
	"github.com/usace/wat-api/wat"
)

type WatHandler struct {
	store   *filestore.FileStore
	queue   *sqs.SQS
	AppPort string
}

func CreateWatHandler() (*WatHandler, error) {
	loader, err := utils.InitLoader("WAT_API")
	wh := WatHandler{}
	store, err := loader.InitStore()
	if err != nil {
		return &wh, err
	}
	wh.store = &store
	sqs, err := loader.InitQueue()
	if err != nil {
		return &wh, err
	}
	wh.queue = sqs
	wh.AppPort = loader.AppPort()
	return &wh, nil
}

const version = "2.0.1 Development"

func (wh *WatHandler) Version(c echo.Context) error {
	return c.String(http.StatusOK, fmt.Sprintf("WAT API Version %s", version))
}
func (wh *WatHandler) Plugins(c echo.Context) error {
	//ping the network to figure out what plugins are active?
	plugins := make([]wat.Plugin, 2)
	plugins[0] = wat.Plugin{Name: "plugin a"}
	plugins[1] = wat.Plugin{Name: "plugin b"}
	return c.JSON(http.StatusOK, plugins)
}
func (wh *WatHandler) ExecuteJob(c echo.Context) error {
	tw := wat.TimeWindow{StartTime: time.Date(2018, 1, 1, 1, 1, 1, 1, time.Local), EndTime: time.Date(2020, time.December, 31, 1, 1, 1, 1, time.Local)}
	sj := wat.StochasticJob{

		TimeWindow:                   tw,
		TotalRealizations:            2,
		EventsPerRealization:         10,
		InitialRealizationSeed:       1234,
		InitialEventSeed:             1234,
		Outputdestination:            "testing",
		Inputsource:                  "testSettings.InputDataDir",
		DeleteOutputAfterRealization: false,
	}
	configs, err := sj.GeneratePayloads(wh.queue)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, configs)
}
