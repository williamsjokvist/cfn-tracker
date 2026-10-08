package config

type BuildConfig struct {
	AppVersion        string `default:"1.0.0"`
	Headless          bool   `envconfig:"HEADLESS" default:"true"`
	BrowserSourcePort int    `envconfig:"BROWSER_SOURCE_PORT" default:"4242"`
}
