package config

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/spf13/viper"
)

type WeaveConfig struct {
	Add_all   bool
	Format    string
	Title     string
	Symbols   map[string]string
	Separator string
}

func MakePresets() {
	viper.SetDefault("add_all", false)
	viper.SetDefault("format", "<type> <symbol>")
	viper.SetDefault("separator", ": ")
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
		log.Fatal(err)
	}

	viper.AddConfigPath(".")
	viper.SetConfigName("weave")
	viper.SetConfigType("toml")
	viper.MergeInConfig()
}

func MakeConfig() WeaveConfig {
	return WeaveConfig{
		Add_all:   viper.GetBool("add_all"),
		Title:     viper.GetString("title"),
		Format:    viper.GetString("format"),
		Symbols:   viper.GetStringMapString("symbols"),
		Separator: viper.GetString("separator"),
	}
}

func (c *WeaveConfig) SymbolChoices(format string) []string {
	choices := make([]string, len(c.Symbols))
	i := 0
	for kind, symbol := range c.Symbols {
		choices[i] = strings.ReplaceAll(strings.ReplaceAll(format, "<type>", kind), "<symbol>", symbol)
		i++
	}
	return choices
}
