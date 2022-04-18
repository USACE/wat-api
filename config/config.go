package config

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
