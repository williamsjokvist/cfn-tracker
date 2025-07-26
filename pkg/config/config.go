package config

type BuildConfig struct {
	AppVersion        string `default:"1.0.0"`
	Headless          bool   `envconfig:"HEADLESS" default:"true"`
	BrowserSourcePort int    `envconfig:"BROWSER_SOURCE_PORT" default:"4242"`
	CapIDEmail        string `envconfig:"CAP_ID_EMAIL" required:"true"`
	CapIDPassword     string `envconfig:"CAP_ID_PASSWORD" required:"true"`
}
