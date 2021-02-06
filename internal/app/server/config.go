package server

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Addr string `json:"addr"`
}

func NewConfig() *Config {
	return &Config{
		Addr: ":8080",
	}
}

func ReadConfig(path string, config *Config) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, config); err != nil {
		return err
	}

	return nil
}
