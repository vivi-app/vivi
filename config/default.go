package config

import (
	"github.com/spf13/viper"
	"github.com/vivi-app/vivi/constant"
)

// fields is the config fields with their default values and descriptions
var fields = []*Field{
	// LOGS
	{
		constant.LogsWrite,
		false,
		"Write logs to file",
	},
	{
		constant.LogsLevel,
		"info",
		`Logs level.
Available options are: (from less to most verbose)
panic, fatal, error, warn, info, debug, trace`,
	},
	// END LOGS

	// DOWNLOADER
	{
		constant.DownloaderPath,
		".",
		"Path to the downloader executable",
	},
	// END DOWNLOADER
}

func setDefaults() {
	Default = make(map[string]*Field, len(fields))
	for _, f := range fields {
		Default[f.Key] = f
		viper.SetDefault(f.Key, f.DefaultValue)
		viper.MustBindEnv(f.Key)
	}
}

var Default map[string]*Field
