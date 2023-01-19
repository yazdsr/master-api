package config

type (
	Config struct {
		App  App
		Psql Psql
	}

	App struct {
		Name   string `envconfig:"MASTER_API_APP_NAME"`
		Port   string `envconfig:"MASTER_API_APP_PORT"`
		Secret string `envconfig:"MASTER_API_APP_SECRET"`
	}

	Psql struct {
		Username string `envconfig:"MASTER_API_PSQL_USERNAME"`
		Password string `envconfig:"MASTER_API_PSQL_PASSWORD"`
		DBName   string `envconfig:"MASTER_API_PSQL_DB_NAME"`
		Host     string `envconfig:"MASTER_API_PSQL_HOST"`
		Port     int    `envconfig:"MASTER_API_PSQL_PORT"`
	}
)
