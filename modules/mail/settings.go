package mail

import (
	"github.com/olebedev/config"
	"github.com/wtfutil/wtf/cfg"
	"log"
)

const (
	defaultFocusable = false
	defaultTitle     = "Mail"
)

// Settings defines the configuration properties for this module
type Settings struct {
	common *cfg.Common

	imapSettings    *imapSettings
	defaultPageSize int `help:"The default number of messages to display per page" values:"Numbers greater than 0"`
	// Define your settings attributes here
}

// NewSettingsFromYAML creates a new settings instance from a YAML config block
func NewSettingsFromYAML(name string, ymlConfig *config.Config, globalConfig *config.Config) *Settings {
	imapConfig, err := ymlConfig.Get("imap")

	var imapCfg imapSettings

	if err != nil {
		imapCfg = imapSettings{
			address: "example.com:993",
		}
		log.Fatal(err)
	} else {
		imapCfg = imapSettings{
			address:  imapConfig.UString("address"),
			username: imapConfig.UString("username"),
			password: imapConfig.UString("password"),
		}
	}

	settings := Settings{
		common:          cfg.NewCommonSettingsFromModule(name, defaultTitle, defaultFocusable, ymlConfig, globalConfig),
		imapSettings:    &imapCfg,
		defaultPageSize: ymlConfig.UInt("defaultPageSize", 10),
		// Configure your settings attributes here. See http://github.com/olebedev/config for type details
	}

	return &settings
}
