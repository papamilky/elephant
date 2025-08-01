package main

import (
	"github.com/abenz1267/elephant/internal/common"
)

type Config struct {
	common.Config           `koanf:",squash"`
	LaunchPrefix            string `koanf:"launch_prefix" desc:"overrides the default app2unit or uwsm prefix, if set. 'CLEAR' to not prefix." default:""`
	Locale                  string `koanf:"locale" desc:"to override systems locale" default:""`
	ShowActions             bool   `koanf:"show_actions" desc:"include application actions, f.e. 'New Private Window' for Firefox" default:"false"`
	ShowGeneric             bool   `koanf:"show_generic" desc:"include generic info when show_actions is true" default:"false"`
	ShowActionsWithoutQuery bool   `koanf:"show_actions_without_query" desc:"show application actions, if the search query is empty" default:"false"`
}

var config *Config

func loadConfig() {
	config = &Config{
		Config:                  common.Config{},
		ShowActions:             false,
		ShowGeneric:             false,
		ShowActionsWithoutQuery: false,
	}

	common.LoadConfig(Name, config)
}
