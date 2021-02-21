// SPDX-License-Identifier: Apache-2.0
// Copyright 2018,2019 Marcus Soll
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	  http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package helper

import (
	"encoding/json"
	"os"
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

	configData, err := os.ReadFile(path.Join(u.HomeDir, ".config/gogitui/config.json"))
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
