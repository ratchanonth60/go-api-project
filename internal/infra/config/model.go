package config

type AppConfig struct {
	Server struct {
		Host string `yaml:"host" env:"HOST" envDefault:"localhost"`
		Port string `yaml:"port" env:"PORT" envDefault:"8000"`
	} `yaml:"server"`
	Database struct {
		Host     string `yaml:"host" env:"POSTGRES_HOST"`
		Port     string `yaml:"port" env:"POSTGRES_PORT" envDefault:"5432"`
		User     string `yaml:"user" env:"POSTGRES_USER"`
		Password string `yaml:"password" env:"POSTGRES_PASSWORD"`
		DBName   string `yaml:"db_name" env:"POSTGRES_DB"`
		SSLMode  string `yaml:"sslmode" env:"POSTGRES_SSLMODE" envDefault:"require"`
	} `yaml:"database"`
	JWT struct {
		Signed string `yaml:"signed" env:"JWT_SIGNED"`
	} `yaml:"jwt"`
	S3 struct {
		Region   string `yaml:"region" env:"AWS_REGION"`
		Bucket   string `yaml:"bucket" env:"AWS_BUCKET"`
		Endpoint string `yaml:"endpoint" env:"AWS_ENDPOINT"`
	} `yaml:"s3"`
	Credentials struct {
		AccessKey string `yaml:"access_key" env:"AWS_ACCESS_KEY_ID"`
		SecretKey string `yaml:"secret_key" env:"AWS_SECRET_ACCESS_KEY"`
	} `yaml:"credentials"`
	Worker struct {
		Broker  string `yaml:"broker" env:"BROKER"`
		Backend string `yaml:"backend" env:"RESULT_BACKEND"`
	} `yaml:"worker"`
	SES struct {
		Region    string `yaml:"region" env:"AWS_REGION_SES"`
		From      string `yaml:"from" env:"EMAIL_FROM"`
		AccessKey string `yaml:"access_key_ses" env:"ACCESS_SES"`
		SecretKey string `yaml:"secret_key_ses" env:"SECRET_SES"`
		Endpoint  string `yaml:"endpoint_ses" env:"ENDPOINT_SES"`
	} `yaml:"ses"`
	SQS struct {
		Region    string `yaml:"region" env:"AWS_REGION_SQS"`
		AccessKey string `yaml:"access_key_sqs" env:"ACCESS_SQS"`
		SecretKey string `yaml:"secret_key_sqs" env:"SECRET_SQS"`
		Endpoint  string `yaml:"endpoint_sqs" env:"ENDPOINT_SQS"`
	} `yaml:"sqs"`
	Redis struct {
		Endpoint string `yaml:"endpoint" env:"REDIS_ENDPOINT"`
		Password string `yaml:"password" env:"REDIS_PASSWORD"`
	} `yaml:"redis"`
}
