package config

type BuildConfig struct {
	AppVersion        string `default:"1.0.0"`
	BrowserSourcePort int    `envconfig:"BROWSER_SOURCE_PORT" default:"4242"`
}
