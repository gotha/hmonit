package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"

	"github.com/Financial-Times/gourmet/config"
	"github.com/Financial-Times/gourmet/log"
)

type Service struct {
	ID   string   `json:"id"`
	Name string   `json:"name"`
	URL  string   `json:"url"`
	Tags []string `json:"tags"`
}

func (s *Service) GetID() string {
	if s.ID == "" {
		algorithm := md5.New()
		_, _ = algorithm.Write([]byte(fmt.Sprintf("%s-%s", s.Name, s.URL)))
		s.ID = hex.EncodeToString(algorithm.Sum(nil))
	}
	return s.ID
}

type appConfig struct {
	SystemCode string `conf:"APP_SYSTEM_CODE" required:"true"`
	LogLevel   string `conf:"LOG_LEVEL" default:"INFO"`
	Server     struct {
		Port         int `conf:"SERVER_PORT" default:"8080"`
		ReadTimeout  int `conf:"SERVER_READ_TIMEOUT" default:"10"`
		WriteTimeout int `conf:"SERVER_WRITE_TIMETOUT" default:"15"`
		IdleTimeout  int `conf:"SERVER_IDLE_TIMEOUT" default:"20"`
	}
	RefreshInterval    int    `conf:"REFRESH_INTERNAL" default:"60"`
	ServicesConfigFile string `conf:"SERVICES_CONFIG_FILE" required:"true"`
	Services           []Service
}

func (c *appConfig) GetLogLevel() log.Level {
	levelMap := map[string]log.Level{
		"TRACE":   log.TraceLevel,
		"DEBUG":   log.DebugLevel,
		"INFO":    log.InfoLevel,
		"WARN":    log.WarnLevel,
		"WARNING": log.WarnLevel,
		"ERROR":   log.ErrorLevel,
	}
	val, exists := levelMap[c.LogLevel]
	if exists {
		return val
	}
	return log.InfoLevel
}

func getConfig() appConfig {
	var err error
	confLoader := config.NewEnvConfigLoader()
	conf := appConfig{}
	err = confLoader.Load(&conf)
	if err != nil {
		panic(fmt.Sprintf("unable to load confiuration: %s", err.Error()))
	}

	dat, err := os.ReadFile(conf.ServicesConfigFile)
	if err != nil {
		panic(fmt.Sprintf("could not read services config file: %s", err.Error()))
	}

	err = json.Unmarshal(dat, &conf.Services)
	if err != nil {
		panic(fmt.Sprintf("could not parse services config file: %s", err.Error()))
	}

	return conf
}
