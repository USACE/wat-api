package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/USACE/filestore"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/labstack/echo/v4"
	"github.com/usace/wat-api/config"
	"github.com/usace/wat-api/wat"
)

type WatHandler struct {
	store *filestore.FileStore
	queue *sqs.SQS
}

func CreateWatHandler(cfg config.WatConfig) (*WatHandler, error) {
	wh := WatHandler{}
	store, err := LoadFileStore(cfg)
	if err != nil {
		return &wh, err
	}
	wh.store = store
	sqs, err := LoadSQS(cfg)
	if err != nil {
		return &wh, err
	}
	wh.queue = sqs
	return &wh, nil
}
func LoadFileStore(cfg config.WatConfig) (*filestore.FileStore, error) {
	s3Conf := filestore.S3FSConfig{
		S3Id:     cfg.AWS_ACCESS_KEY_ID,
		S3Key:    cfg.AWS_SECRET_ACCESS_KEY,
		S3Region: cfg.AWS_DEFAULT_REGION,
		S3Bucket: cfg.S3_BUCKET,
	}
	if cfg.S3_MOCK {
		s3Conf.Mock = cfg.S3_MOCK
		s3Conf.S3DisableSSL = cfg.S3_DISABLE_SSL
		s3Conf.S3ForcePathStyle = cfg.S3_FORCE_PATH_STYLE
		s3Conf.S3Endpoint = cfg.S3_ENDPOINT
	}
	fmt.Println(s3Conf)

	fs, err := filestore.NewFileStore(s3Conf)

	if err != nil {
		return nil, err
	}

	return &fs, nil
}
func LoadSQS(cfg config.WatConfig) (*sqs.SQS, error) {
	creds := credentials.NewStaticCredentials(cfg.AWS_ACCESS_KEY_ID, cfg.AWS_SECRET_ACCESS_KEY, "")
	awscfg := aws.NewConfig().WithRegion(cfg.AWS_DEFAULT_REGION).WithCredentials(creds)
	sess, err := session.NewSession(awscfg)
	if err != nil {
		return nil, err
	}

	sqs := sqs.New(sess, aws.NewConfig().WithEndpoint(cfg.SQS_ENDPOINT))
	return sqs, nil
}
func LoadRedis() {

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
