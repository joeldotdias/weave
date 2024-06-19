package config

import (
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/viper"
)

type WeaveConfig struct {
	Add_all bool
	Title   string
	Symbols map[string]string
}

func MakePresets() {
	viper.SetDefault("add_all", false)
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	if runtime.GOOS == "windows" {
		cfg_dir := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if cfg_dir == "" {
			cfg_dir = os.Getenv("USERPROFILE")
		}
		viper.AddConfigPath(cfg_dir + "\\weave")
	} else {
		viper.AddConfigPath("$XDG_CONFIG_HOME/weave")
		viper.AddConfigPath("$HOME/.config/weave")
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("No config file")
		}
	}

	viper.AddConfigPath(".")
	viper.SetConfigName("weave")
	viper.SetConfigType("toml")
	if err := viper.MergeInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("No local env found. Using global config")
		}
	}
}

func MakeConfig() WeaveConfig {
	return WeaveConfig{
		Add_all: viper.GetBool("add_all"),
		Title:   viper.GetString("title"),
		Symbols: viper.GetStringMapString("symbols"),
	}
}

func (c *WeaveConfig) SymbolChoices() []string {
	choices := make([]string, len(c.Symbols))
	i := 0
	for name, emoji := range c.Symbols {
		choices[i] = fmt.Sprintf("%s %s", name, emoji)
		i++
	}
	return choices
}
