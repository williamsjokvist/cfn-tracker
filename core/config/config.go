package config

type Config struct {
	AppVersion        string `envconfig:"APP_VERSION" default:"0.0.0"`
	Headless          string `envconfig:"HEADLESS" default:"true"`
	BrowserSourcePort string `envconfig:"BROWSER_SOURCE_PORT" default:"4242"`
	SteamUsername     string `envconfig:"STEAM_USERNAME" required:"true"`
	SteamPassword     string `envconfig:"STEAM_PASSWORD" required:"true"`
	CapIDEmail        string `envconfig:"CAP_ID_EMAIL" required:"true"`
	CapIDPassword     string `envconfig:"CAP_ID_PASSWORD" required:"true"`
}
