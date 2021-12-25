package repo

import "github.com/kelseyhightower/envconfig"

type (
	DBParameters struct {
		Host      string `envconfig:"POSTGRES_HOST" default:"localhost"`
		Port      int    `envconfig:"POSTGRES_PORT" default:"5432"`
		User      string `envconfig:"POSTGRES_USER" default:"mrpostalofficer"`
		Pass      string `envconfig:"POSTGRES_PASSWORD"`
		Name      string `envconfig:"POSTGRES_NAME" default:"packagetracer"`
		Options   string `envconfig:"POSTGRES_OPTIONS" default:"sslmode=disable"`
		TableName string `envconfig:"TABLE_NAME" default:"customs"`
	}
)

var (
	DBConfig DBParameters
)

func LoadDBConfig() {
	envconfig.MustProcess("", &DBConfig)
}
