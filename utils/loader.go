package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/USACE/filestore"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/batch"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/go-redis/redis"
	"github.com/kelseyhightower/envconfig"
	"github.com/usace/wat-api/config"
	"github.com/usace/wat-api/model"
	"gopkg.in/yaml.v2"
)

type ServicesLoader struct {
	config config.WatConfig
}

func InitLoaderWithConfig(prefix string, cfg config.WatConfig) (ServicesLoader, error) {
	ldr := ServicesLoader{}
	ldr.config = cfg
	return ldr, nil
}
func InitLoader(prefix string) (ServicesLoader, error) {
	var cfg config.WatConfig
	ldr := ServicesLoader{}
	if err := envconfig.Process(prefix, &cfg); err != nil {
		return ldr, err
	}
	ldr.config = cfg
	return ldr, nil
}
func (sl ServicesLoader) AppPort() string {
	return sl.config.APP_PORT
}
func (sl ServicesLoader) EnvironmentVariables() []string {
	return sl.config.EnvironmentVariables()
}
func (sl ServicesLoader) Config() config.WatConfig {
	return sl.config
}
func (sl ServicesLoader) InitQueue() (*sqs.SQS, error) {
	creds := credentials.NewStaticCredentials(
		sl.config.AWS_ACCESS_KEY_ID,
		sl.config.AWS_SECRET_ACCESS_KEY,
		"",
	)
	awscfg := aws.NewConfig().WithRegion(sl.config.AWS_DEFAULT_REGION).WithCredentials(creds)
	sess, err := session.NewSession(awscfg)
	if err != nil {
		return nil, err
	}
	sqs := sqs.New(sess, aws.NewConfig().WithEndpoint(sl.config.SQS_ENDPOINT))
	return sqs, nil
}
func (sl ServicesLoader) InitStore() (filestore.FileStore, error) {
	//initalize S3 Store
	mock := sl.config.S3_MOCK
	s3Conf := filestore.S3FSConfig{
		S3Id:     sl.config.AWS_ACCESS_KEY_ID,
		S3Key:    sl.config.AWS_SECRET_ACCESS_KEY,
		S3Region: sl.config.AWS_DEFAULT_REGION,
		S3Bucket: sl.config.S3_BUCKET,
	}
	if mock {
		s3Conf.Mock = mock
		s3Conf.S3DisableSSL = sl.config.S3_DISABLE_SSL
		s3Conf.S3ForcePathStyle = sl.config.S3_FORCE_PATH_STYLE
		s3Conf.S3Endpoint = sl.config.S3_ENDPOINT
	}
	fmt.Println(s3Conf)

	fs, err := filestore.NewFileStore(s3Conf)

	if err != nil {
		log.Fatal(err)
	}

	return fs, nil
}
func (sl ServicesLoader) InitRedis() (*redis.Client, error) {

	client := redis.NewClient(&redis.Options{
		Addr:     sl.config.REDIS_HOST + ":" + sl.config.REDIS_PORT,
		Password: sl.config.REDIS_PASSWORD,
		DB:       0,
	})
	return client, nil
}
func (sl ServicesLoader) InitBatch() (*batch.Batch, error) {
	/*sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentialsFromCreds(credentials.Value{
		   AccessKeyID:     accessKeyValue,
		   SecretAccessKey: secret,
		}),
		Region: aws.String("us-east-1")},
	 )*/
	creds := credentials.NewStaticCredentials(
		sl.config.AWS_ACCESS_KEY_ID,
		sl.config.AWS_SECRET_ACCESS_KEY,
		"",
	)
	awscfg := aws.NewConfig().WithRegion(sl.config.AWS_DEFAULT_REGION).WithCredentials(creds)
	sess, err := session.NewSession(awscfg)
	if err != nil {
		return nil, err
	}
	batchClient := batch.New(sess)
	return batchClient, nil
}
func LoadJsonPluginModelFromS3(filepath string, fs filestore.FileStore, pluginModel interface{}) error {
	fmt.Println("reading:", filepath)
	data, err := fs.GetObject(filepath)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(data)
	if err != nil {
		return err
	}

	// fmt.Println("read:", string(body))
	errjson := json.Unmarshal(body, &pluginModel)
	if errjson != nil {
		fmt.Println("error:", errjson)
		return errjson
	}

	return nil

}

// LoadModelPayload
func LoadModelPayloadFromS3(payloadFile string, fs filestore.FileStore) (model.ModelPayload, error) {
	var p model.ModelPayload
	fmt.Println("reading payload:", payloadFile)
	data, err := fs.GetObject(payloadFile)
	if err != nil {
		return p, err
	}

	body, err := ioutil.ReadAll(data)
	if err != nil {
		return p, err
	}
	//fmt.Println(string(body))
	err = yaml.Unmarshal(body, &p)
	if err != nil {
		return p, err
	}
	//fmt.Println(p)
	return p, nil
}

// UpLoadToS3
func UpLoadToS3(newS3Path string, fileBytes []byte, fs filestore.FileStore) (filestore.FileOperationOutput, error) {
	var repsonse *filestore.FileOperationOutput
	var err error
	repsonse, err = fs.PutObject(newS3Path, fileBytes)
	if err != nil {
		return *repsonse, err
	}

	return *repsonse, err
}
