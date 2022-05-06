package handler

import (
	"fmt"
	"net/http"

	"github.com/USACE/filestore"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/go-redis/redis"
	"github.com/labstack/echo/v4"
	"github.com/usace/wat-api/config"
	"github.com/usace/wat-api/utils"
)

type WatHandler struct {
	store   filestore.FileStore
	queue   *sqs.SQS
	cache   *redis.Client
	AppPort string
	config  config.WatConfig
}

func CreateWatHandler() (*WatHandler, error) {
	loader, err := utils.InitLoader("WAT_API")

	wh := WatHandler{}
	store, err := loader.InitStore()
	if err != nil {
		return &wh, err
	}
	wh.store = store
	sqs, err := loader.InitQueue()
	if err != nil {
		return &wh, err
	}
	wh.queue = sqs
	cache, err := loader.InitRedis()
	if err != nil {
		return &wh, err
	}
	wh.cache = cache
	wh.AppPort = loader.AppPort()
	wh.config = loader.Config()
	return &wh, nil
}

const version = "2.0.1 Development"

func (wh *WatHandler) Version(c echo.Context) error {
	return c.String(http.StatusOK, fmt.Sprintf("WAT API Version %s", version))
}
func (wh *WatHandler) Plugins(c echo.Context) error {
	//ping the network to figure out what plugins are active?
	plugins := MockPlugins()
	return c.JSON(http.StatusOK, plugins)
}
func (wh *WatHandler) ExecuteJob(c echo.Context) error {
	sj := MockStochasticJob(wh.config)
	err := sj.GeneratePayloads(wh.queue, wh.store, wh.cache, wh.config)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.String(http.StatusOK, "Compute Started")
}
