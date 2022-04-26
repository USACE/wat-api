package config

import "fmt"

type WatConfig struct {
	APP_PORT              string
	SKIP_JWT              bool
	AWS_ACCESS_KEY_ID     string
	AWS_SECRET_ACCESS_KEY string
	AWS_DEFAULT_REGION    string
	AWS_S3_REGION         string
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
	ret := make([]string, 15)
	ret[0] = "APP_PORT=" + wc.APP_PORT
	ret[1] = fmt.Sprintf("SKIP_JWT=%v", wc.SKIP_JWT)
	ret[2] = "AWS_ACCESS_KEY_ID=" + wc.AWS_ACCESS_KEY_ID
	ret[3] = "AWS_SECRET_ACCESS_KEY=" + wc.AWS_SECRET_ACCESS_KEY
	ret[4] = "AWS_DEFAULT_REGION=" + wc.AWS_DEFAULT_REGION
	ret[5] = "AWS_S3_REGION=" + wc.AWS_S3_REGION
	ret[6] = fmt.Sprintf("S3_MOCK=%v", wc.S3_MOCK)
	ret[7] = "S3_BUCKET=" + wc.S3_BUCKET
	ret[8] = "S3_ENDPOINT=" + wc.S3_ENDPOINT
	ret[9] = fmt.Sprintf("S3_DISABLE_SSL=%v", wc.S3_DISABLE_SSL)
	ret[10] = fmt.Sprintf("S3_FORCE_PATH_STYLE=%v", wc.S3_FORCE_PATH_STYLE)
	ret[11] = "REDIS_HOST=" + wc.REDIS_HOST
	ret[12] = "REDIS_PORT=" + wc.REDIS_PORT
	ret[13] = "REDIS_PASSWORD=" + wc.REDIS_PASSWORD
	ret[14] = "SQS_ENDPOINT=" + wc.SQS_ENDPOINT
	return ret
}
