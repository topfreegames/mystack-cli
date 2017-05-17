// mystack-cli
// https://github.com/topfreegames/mystack-cli
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2016 Top Free Games <backend@tfgco.com>

package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
)

// Config struct
type Config struct {
	Token          string `json:"token"`
	ControllerURL  string `json:"controllerUrl"`
	ControllerHost string `json:"controllerHost"`
	LoggerHost     string `json:"loggerHost"`
	Env            string `json:"env"`
}

// NewConfig ctor
func NewConfig(env, token, controllerURL string, hosts map[string]string) *Config {
	c := &Config{
		Token:          token,
		ControllerURL:  controllerURL,
		ControllerHost: hosts["controller"],
		LoggerHost:     hosts["logger"],
		Env:            env,
	}
	return c
}

// ReadConfig from file
func ReadConfig(fs FileSystem, env string) (*Config, error) {
	cfgPath, err := getConfigPathForEnv(env)
	if err != nil {
		return nil, err
	}
	if _, err := fs.Stat(cfgPath); fs.IsNotExist(err) {
		return nil, err
	}
	bts, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		return nil, err
	}
	c := &Config{}
	err = json.Unmarshal(bts, c)
	if err != nil {
		return nil, err
	}
	c.Env = env
	return c, nil
}

func getConfigDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	cfgDir := filepath.Join(usr.HomeDir, ".mystack")
	return cfgDir, nil
}

func getConfigPathForEnv(env string) (string, error) {
	cfgDir, err := getConfigDir()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/mystack-%s.json", cfgDir, env), nil
}

// Write the config file to disk
func (c *Config) Write(fs FileSystem) error {
	cfgPath, err := getConfigPathForEnv(c.Env)
	if err != nil {
		return err
	}
	cfg, err := json.Marshal(c)
	if err != nil {
		return err
	}
	cfgDir, err := getConfigDir()
	if err != nil {
		return err
	}
	err = fs.MkdirAll(cfgDir, os.ModePerm)
	if err != nil {
		return err
	}
	file, err := fs.Create(cfgPath)
	if err != nil {
		return err
	}
	_, err = file.Write(cfg)
	if err != nil {
		return err
	}
	return nil
}
