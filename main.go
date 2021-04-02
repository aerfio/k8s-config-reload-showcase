/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	c := New()
	for {
		logrus.Info(c.v.GetString("DUPA"))
		time.Sleep(2 * time.Second)
	}
}

type Configuration struct {
	v *viper.Viper
}

const (
	// Constants for viper variable names. Will be used to set
	// default values as well as to get each value

	varLogLevel = "log.level"
)

func New() *Configuration {
	c := Configuration{
		v: viper.New(),
	}

	c.v.SetDefault(varLogLevel, "info")
	c.v.AutomaticEnv()
	// c.v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	c.v.SetTypeByDefaultValue(true)
	c.v.SetConfigFile("/data/config.yaml")
	err := c.v.ReadInConfig() // Find and read the config file
	logrus.Warn("loading config")
	// just use the default value(s) if the config file was not found
	if _, ok := err.(*os.PathError); ok {
		logrus.Warn("no config file not found. Using default values")
	} else if err != nil { // Handle other errors that occurred while reading the config file
		panic(fmt.Errorf("fatal error while reading the config file: %s", err))
	}

	// monitor the changes in the config file
	c.v.WatchConfig()
	c.v.OnConfigChange(func(e fsnotify.Event) {
		logrus.WithField("file", e.Name).WithField("eventName", e.String()).Warn("Config file changed")
	})
	return &c
}

// GetLogLevel returns the log level
func (c *Configuration) GetLogLevel() string {
	return c.v.GetString(varLogLevel)
}
