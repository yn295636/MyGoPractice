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
	MongoAddr  string   `toml:"mongo_addr"`
	RedisAddr  string   `toml:"redis_addr"`
	DbAddr     string   `toml:"db_addr"`
	DbUser     string   `toml:"db_user"`
	DbPassword string   `toml:"db_password"`
	EtcdAddrs  []string `toml:"etcd_addrs"`
}

func SetSettings(s *Settings) {
	settings = s
}
func GetSettings() *Settings {
	if settings == nil {
		settings = &Settings{
			ListenPort: ":50051",
			MongoAddr:  ":27017",
			RedisAddr:  "127.0.0.1:6379",
			DbAddr:     "127.0.0.1:3306",
			DbUser:     "root",
			DbPassword: "Mygopractice123!",
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
	flagSet := flag.NewFlagSet("greeter_service", flag.ExitOnError)
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
