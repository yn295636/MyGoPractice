package main

import (
	"flag"
	"github.com/BurntSushi/toml"
	"log"
	"os"
)

var (
	settings *Settings
)

type Config map[string]interface{}

func (t *Config) Validate() {
}

type Settings struct {
	ListenPort string   `toml:"listen_port"`
	EtcdAddrs  []string `toml:"etcd_addrs"`
}

func SetSettings(s *Settings) {
	settings = s
}
func GetSettings() *Settings {
	if settings == nil {
		settings = &Settings{
			ListenPort: "8081",
			EtcdAddrs: []string{
				"localhost:2379",
				"localhost:22379",
				"localhost:32379",
			},
		}
	}
	return settings
}

func LoadConfig() {
	setting := GetSettings()
	flagSet := flag.NewFlagSet("apigateway", flag.ExitOnError)
	flagSet.String("config", "", "config file path")
	flagSet.Parse(os.Args[1:])

	configFile := flagSet.Lookup("config").Value.String()
	if configFile != "" {
		_, err := toml.DecodeFile(configFile, &setting)
		if err != nil {
			log.Printf("Failed to load config file=%s, err=%s", configFile, err.Error())
		}
	} else {
		log.Println("Warning: config file doesn't exist.")
	}
	log.Println("setting: ", setting)
}
