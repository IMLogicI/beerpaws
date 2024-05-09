package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

var (
	Token  string
	Prefix string
	config *configStruct
	DBConf DBConfig
)

type DBConfig struct {
	Host     string
	Port     int64
	User     string
	Password string
	DBName   string
}

type configStruct struct {
	Token    string `json:"token"`
	Prefix   string `json:"prefix"`
	Host     string `json:"host"`
	Port     int64  `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
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
	DBConf = DBConfig{
		Host:     config.Host,
		Port:     config.Port,
		User:     config.User,
		Password: config.Password,
		DBName:   config.DBName,
	}
	return nil
}
