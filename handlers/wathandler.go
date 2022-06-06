package handler

import (
	"fmt"
	"net/http"

	"github.com/USACE/filestore"
	"github.com/aws/aws-sdk-go/service/batch"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/go-redis/redis"
	"github.com/labstack/echo/v4"
	"github.com/usace/wat-api/config"
	"github.com/usace/wat-api/model"
	"github.com/usace/wat-api/utils"
	"github.com/usace/wat-api/wat"
)

type WatHandler struct {
	store         filestore.FileStore
	queue         *sqs.SQS
	cache         *redis.Client
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
	plugins := make([]model.Plugin, 0)
	dag := MockDag()
	for _, m := range dag.Nodes {
		plugins = append(plugins, m.Plugin)
	}
	return c.JSON(http.StatusOK, plugins)
}
func (wh *WatHandler) ExecuteJob(c echo.Context) error {
	sj := wat.StochasticJob{}
	if err := c.Bind(&sj); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	err := sj.GeneratePayloads(wh.queue, wh.store, wh.cache, wh.config, wh.captainCrunch)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.String(http.StatusOK, "Compute Started")
}
