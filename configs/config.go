package configs

import (
	"flag"
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Data  []int   `toml:"data"`
	Gamma float64 `toml:"gamma"`
	Time1 int     `toml:"time1"`
	Time2 int     `toml:"time2"`
}

var (
	parsedConfigPath string
)

func init() {
	flag.StringVar(&parsedConfigPath, "config-path", "config.toml", "path to config file")
}

func NewConfig() (c Config) {
	flag.Parse()

	_, err := toml.DecodeFile(parsedConfigPath, &c)
	if err != nil {
		log.Println(`Broken or missing config file. 
Provide path to your config file by using "-config-path=<path_to_your_config_file>.toml" or put it in the project's root`)
		os.Exit(1)
	}

	return
}
