package beacons

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Pid    string
	Script string
	Stream map[string]map[string]interface{}
	Agent  map[string]map[string]interface{}

	filename string
}

func LoadConfig(filename string, config *Config) error {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(content, &config); err != nil {
		return err
	}
	config.filename = filename
	return nil
}

func (config Config) SaveTo(filename string) error {
	content, err := json.Marshal(config)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(filename, content, 0600); err != nil {
		return err
	}
	return nil
}

func (config Config) Save() error {
	return config.SaveTo(config.filename)
}
