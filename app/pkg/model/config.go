package model

type RuntimeConfig struct {
	GUI  GUIConfig  `ini:"gui"`
	Data DataConfig `ini:"data"`
}

type GUIConfig struct {
	Locale            string    `ini:"locale" json:"locale" default:"en-GB"`
	Theme             ThemeName `ini:"theme" json:"theme" default:"tekken"`
	SideBar           bool      `ini:"sidebar" json:"sidebar" default:"true"`
	LogFile           bool      `ini:"log_file" json:"logFile" default:"true"`
	BrowserSourcePort uint16    `ini:"browser_source_port" json:"browserSourcePort" default:"4242"`
}

type DataConfig struct {
	// Empty means the default app data directory.
	Dir string `ini:"dir"`
}

type DataLocation struct {
	Dir        string `json:"dir"`
	DefaultDir string `json:"defaultDir"`
	// The database is only swapped on start, so a move or an import waits for a restart.
	RestartRequired bool `json:"restartRequired"`
}
