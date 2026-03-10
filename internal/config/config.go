package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type Environment struct {
	AppEnv                    string `envconfig:"APP_ENV" default:"local"`
	AppPort                   string `envconfig:"APP_PORT" default:"8080"`
	SessionSecret             string `envconfig:"SESSION_SECRET"`
	GoogleClientID            string `envconfig:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret        string `envconfig:"GOOGLE_CLIENT_SECRET"`
	UserGoogleRedirectURL     string `envconfig:"USER_GOOGLE_REDIRECT_URL"`
	OperatorGoogleRedirectURL string `envconfig:"OPERATOR_GOOGLE_REDIRECT_URL"`
	GA4GoogleRedirectURL      string `envconfig:"GA4_GOOGLE_REDIRECT_URL"`
	GSCGoogleRedirectURL      string `envconfig:"GSC_GOOGLE_REDIRECT_URL"`
	GBPGoogleRedirectURL      string `envconfig:"GBP_GOOGLE_REDIRECT_URL"`
	MetaAppID                 string `envconfig:"META_APP_ID"`
	MetaAppSecret             string `envconfig:"META_APP_SECRET"`
	InstagramRedirectURL      string `envconfig:"INSTAGRAM_REDIRECT_URL"`
	DBHost                    string `envconfig:"DB_HOST" default:"localhost"`
	DBPort                    string `envconfig:"DB_PORT" default:"3306"`
	DBUser                    string `envconfig:"DB_USER"`
	DBPassword                string `envconfig:"DB_PASSWORD"`
	DBName                    string `envconfig:"DB_NAME"`
	FrontendURL               string `envconfig:"FRONTEND_URL" default:"http://localhost:5173"`
}

var Env Environment

func init() {
	err := envconfig.Process("", &Env)
	if err != nil {
		log.Fatal(err)
	}
}
