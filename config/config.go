package config

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/batch"
)

type WatConfig struct {
	APP_PORT              string
	SKIP_JWT              bool
	AWS_ACCESS_KEY_ID     string
	AWS_SECRET_ACCESS_KEY string
	AWS_DEFAULT_REGION    string
	AWS_S3_REGION         string
	AWS_REGION            string
	AWS_BUCKET            string
	S3_MOCK               bool
	S3_BUCKET             string
	S3_ENDPOINT           string
	S3_DISABLE_SSL        bool
	S3_FORCE_PATH_STYLE   bool
	REDIS_HOST            string
	REDIS_PORT            string
	REDIS_PASSWORD        string
	SQS_ENDPOINT          string
}

func (wc WatConfig) EnvironmentVariables() []string {
	ret := make([]string, 17)
	ret[0] = "APP_PORT=" + wc.APP_PORT
	ret[1] = fmt.Sprintf("SKIP_JWT=%v", wc.SKIP_JWT)
	ret[2] = "AWS_ACCESS_KEY_ID=" + wc.AWS_ACCESS_KEY_ID
	ret[3] = "AWS_SECRET_ACCESS_KEY=" + wc.AWS_SECRET_ACCESS_KEY
	ret[4] = "AWS_DEFAULT_REGION=" + wc.AWS_DEFAULT_REGION
	ret[5] = "AWS_S3_REGION=" + wc.AWS_S3_REGION
	ret[6] = "AWS_REGION=" + wc.AWS_REGION
	ret[7] = "AWS_BUCKET=" + wc.AWS_BUCKET
	ret[8] = fmt.Sprintf("S3_MOCK=%v", wc.S3_MOCK)
	ret[9] = "S3_BUCKET=" + wc.S3_BUCKET
	ret[10] = "S3_ENDPOINT=" + wc.S3_ENDPOINT
	ret[11] = fmt.Sprintf("S3_DISABLE_SSL=%v", wc.S3_DISABLE_SSL)
	ret[12] = fmt.Sprintf("S3_FORCE_PATH_STYLE=%v", wc.S3_FORCE_PATH_STYLE)
	ret[13] = "REDIS_HOST=" + wc.REDIS_HOST
	ret[14] = "REDIS_PORT=" + wc.REDIS_PORT
	ret[15] = "REDIS_PASSWORD=" + wc.REDIS_PASSWORD
	ret[16] = "SQS_ENDPOINT=" + wc.SQS_ENDPOINT
	return ret
}
func toBatchKeyValuePair(key string, value string) *batch.KeyValuePair {
	keyvalue := batch.KeyValuePair{
		Name:  aws.String(key),
		Value: aws.String(value),
	}
	return &keyvalue
}
func (wc WatConfig) BatchEnvironmentVariables() []*batch.KeyValuePair {
	ret := make([]*batch.KeyValuePair, 13)
	ret[0] = toBatchKeyValuePair("AWS_ACCESS_KEY_ID", wc.AWS_ACCESS_KEY_ID)
	ret[1] = toBatchKeyValuePair("AWS_SECRET_ACCESS_KEY", wc.AWS_SECRET_ACCESS_KEY)
	ret[2] = toBatchKeyValuePair("AWS_DEFAULT_REGION", wc.AWS_DEFAULT_REGION)
	ret[3] = toBatchKeyValuePair("AWS_S3_REGION", wc.AWS_S3_REGION)
	ret[4] = toBatchKeyValuePair("S3_MOCK", fmt.Sprintf("%v", wc.S3_MOCK))
	ret[5] = toBatchKeyValuePair("S3_BUCKET", wc.S3_BUCKET)
	ret[6] = toBatchKeyValuePair("S3_ENDPOINT", wc.S3_ENDPOINT)
	ret[7] = toBatchKeyValuePair("S3_DISABLE_SSL", fmt.Sprintf("%v", wc.S3_DISABLE_SSL))
	ret[8] = toBatchKeyValuePair("S3_FORCE_PATH_STYLE", fmt.Sprintf("%v", wc.S3_FORCE_PATH_STYLE))
	ret[9] = toBatchKeyValuePair("REDIS_HOST", wc.REDIS_HOST)
	ret[10] = toBatchKeyValuePair("REDIS_PORT", wc.REDIS_PORT)
	ret[11] = toBatchKeyValuePair("REDIS_PASSWORD", wc.REDIS_PASSWORD)
	ret[12] = toBatchKeyValuePair("SQS_ENDPOINT", wc.SQS_ENDPOINT)
	return ret
}
