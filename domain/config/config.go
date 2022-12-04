package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	App      App      `required:"true" envconfig:"APP"`
	DB       DB       `required:"true" envconfig:"DB"`
	Firebase Firebase `required:"true" envconfig:"FIREBASE"`
	Sendgrid Sendgrid `required:"true" envconfig:"SENDGRID"`
	Slack    Slack    `required:"true" envconfig:"SLACK"`
	Google   Google   `required:"true" envconfig:"GOOGLE"`
	Line     Line     `required:"true" envconfig:"LINE"`
}

func New() (Config, error) {
	var c Config
	err := envconfig.Process("", &c)
	if err != nil {
		return Config{}, err
	}
	return c, nil
}

type App struct {
	Env            string   `required:"true" split_words:"true"`
	Service        string   `required:"true" split_words:"true"`
	BatchType      string   `required:"true" split_words:"true"`
	Port           int      `required:"true" split_words:"true"`
	BasicUsers     []string `required:"true" split_words:"true"`
	BasicPasswords []string `required:"true" split_words:"true"`
	CorsDomains    []string `required:"true" split_words:"true"`
}

type DB struct {
	Host string `required:"true" split_words:"true"`
	Port int    `required:"true" split_words:"true"`
	Name string `required:"true" split_words:"true"`
	User string `required:"true" split_words:"true"`
	Pass string `required:"true" split_words:"true"`
}

type Firebase struct {
	JSONFilePath string `required:"true" split_words:"true"`
	WebAPIKey    string `required:"true" split_words:"true"`
}

type Sendgrid struct {
	APIKey string `required:"true" split_words:"true"`
}

type Slack struct {
	AccessToken string `required:"true" split_words:"true"`
	ChanelID    string `required:"true" split_words:"true"`
}

type Google struct {
	ApplicationCredentials string `required:"true" split_words:"true"`
}

type Line struct {
	ChannelSecret      string `required:"true" split_words:"true"`
	ChannelAccessToken string `required:"true" split_words:"true"`
}
