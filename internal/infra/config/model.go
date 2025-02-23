package config

type AppConfig struct {
	Server struct {
		Port string `yaml:"port" env:"PORT"`
	} `yaml:"server"`
	Database struct {
		Host     string `yaml:"host" env:"POSTGRES_HOST"`
		Port     string `yaml:"port" env:"POSTGRES_PORT"`
		User     string `yaml:"user" env:"POSTGRES_USER"`
		Password string `yaml:"password" env:"POSTGRES_PASSWORD"`
		DBName   string `yaml:"db_name" env:"POSTGRES_DB"`
	} `yaml:"database"`
	Cache struct {
		Address     string `yaml:"address" env:"REDIS_HOST"`
		Port        string `yaml:"port" env:"REDIS_PORT"`
		MaxSizePool string `yaml:"max_size_pool" env:"REDIS_MAX_POOL"`
	} `yaml:"cache"`
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
}
