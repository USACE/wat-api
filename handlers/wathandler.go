package handler

import (
	"fmt"
	"net/http"

	"github.com/USACE/filestore"
	"github.com/aws/aws-sdk-go/service/batch"
	"github.com/labstack/echo/v4"
	"github.com/usace/wat-api/config"
	"github.com/usace/wat-api/utils"
)

type WatHandler struct {
	store         filestore.FileStore
	captainCrunch *batch.Batch
	AppPort       string
	config        config.WatConfig
}

func CreateWatHandlerFromConfig(config config.WatConfig) (*WatHandler, error) {
	wh := WatHandler{}
	loader, err := utils.InitLoaderWithConfig("WAT_API", config)
	if err != nil {
		return &wh, err
	}
	store, err := loader.InitStore()
	if err != nil {
		return &wh, err
	}
	wh.store = store
	awsBatch, err := loader.InitBatch()
	if err != nil {
		return &wh, err
	}
	wh.captainCrunch = awsBatch
	wh.AppPort = loader.AppPort()
	wh.config = loader.Config()
	return &wh, nil
}
func CreateWatHandler() (*WatHandler, error) {
	wh := WatHandler{}
	loader, err := utils.InitLoader("WAT_API")
	if err != nil {
		return &wh, err
	}
	store, err := loader.InitStore()
	if err != nil {
		return &wh, err
	}
	wh.store = store
	awsBatch, err := loader.InitBatch()
	if err != nil {
		return &wh, err
	}
	wh.captainCrunch = awsBatch
	wh.AppPort = loader.AppPort()
	wh.config = loader.Config()
	return &wh, nil
}

const version = "2.0.1 Development"

func (wh *WatHandler) Version(c echo.Context) error {
	return c.String(http.StatusOK, fmt.Sprintf("WAT API Version %s", version))
}
func (wh WatHandler) Config() config.WatConfig {
	return wh.config
}
func (wh *WatHandler) Plugins(c echo.Context) error {
	//ping the network to figure out what plugins are active?
	plugins := make([]string, 0)
	return c.JSON(http.StatusOK, plugins)
}
func (wh *WatHandler) ExecuteJob(c echo.Context) error {
	fmt.Println("executing job")
	return c.String(http.StatusOK, "Compute Started")
}
