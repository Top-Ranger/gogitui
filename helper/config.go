// SPDX-License-Identifier: MIT

package helper

import (
	"encoding/json"
	"os"
	"io/ioutil"
	"os/user"
	"path"
)

// This struct holds all information
type Config struct {
	// All repositories
	Repositories []string
}

// Loads configuration. Will return an empty configuration if the file does not exists
func LoadConfig() (Config, error) {
	var config Config

	u, err := user.Current()
	if err != nil {
		return config, err
	}

	configData, err := ioutil.ReadFile(path.Join(u.HomeDir, ".config/gogitui/config.json"))
	if os.IsNotExist(err) {
		return config, nil
	}
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(configData, &config)
	return config, err
}

// Saves the configuration to disk
func (config *Config) SaveConfig() error {
	u, err := user.Current()
	if err != nil {
		return err
	}
	err = os.MkdirAll(path.Join(u.HomeDir, ".config/gogitui/"), os.FileMode(0777))
	if err != nil {
		return err
	}

	configData, err := json.Marshal(config)
	if err != nil {
		return err
	}

	file, err := os.Create(path.Join(u.HomeDir, ".config/gogitui/config.json"))
	if err != nil {
		return err
	}

	_, err = file.Write(configData)
	return err
}