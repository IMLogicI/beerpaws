package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

var (
	Token           string
	Prefix          string
	config          *configStruct
	DBConf          DBConfig
	ChannelsID      map[string]interface{}
	ApplicationID   string
	GuildID         string
	ChannelsAdminID string
)

type DBConfig struct {
	Host     string
	Port     int64
	User     string
	Password string
	DBName   string
}

type configStruct struct {
	Token           string   `json:"token"`
	Prefix          string   `json:"prefix"`
	Host            string   `json:"host"`
	Port            int64    `json:"port"`
	User            string   `json:"user"`
	Password        string   `json:"password"`
	DBName          string   `json:"dbname"`
	ChannelsID      []string `json:"channels_id"`
	ApplicationID   string   `json:"application_id"`
	GuildID         string   `json:"guild_id"`
	ChannelsAdminID string   `json:"channels_admin_id"`
}

func ReadConfig() error {
	file, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Fatal(err)
		return err
	}
	err = json.Unmarshal(file, &config)
	if err != nil {
		log.Fatal(err)
		return err
	}
	Token = config.Token
	Prefix = config.Prefix
	ApplicationID = config.ApplicationID
	GuildID = config.GuildID
	ChannelsAdminID = config.ChannelsAdminID
	ChannelsID = make(map[string]interface{})
	for _, channelID := range config.ChannelsID {
		ChannelsID[channelID] = ""
	}

	DBConf = DBConfig{
		Host:     config.Host,
		Port:     config.Port,
		User:     config.User,
		Password: config.Password,
		DBName:   config.DBName,
	}
	return nil
}
