package config

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	"github.com/rwbm/go-tools/files"
	yaml "gopkg.in/yaml.v2"
)

// Configuration is the structure used to hold configuration from config.yml
type Configuration struct {
	Server struct {
		Name         string `yaml:"name"`
		Port         string `yaml:"port"`
		ReadTimeout  int    `yaml:"read_timeout"`
		WriteTimeout int    `yaml:"write_timeout"`
		DryRun       bool   `yaml:"dry_run"`
	} `yaml:"server"`
	Database struct {
		Filename string `yaml:"filename"`
	} `yaml:"database"`
	Template struct {
		Base           string `yaml:"base_location"`
		ProcessedOK    string `yaml:"processed_ok"`
		ProcessedError string `yaml:"processed_error"`
		CheckCycle     int    `yaml:"check_cycle"`
	} `yaml:"template"`
}

// Load reads application settings in the indicated file
func Load(path string) (cfg *Configuration, err error) {
	if files.Exists(path) {
		if cfg, err = loadSettingsFromFile(path); err == nil {
			replaceCustomVars(cfg)
		}
	} else {
		return nil, fmt.Errorf("file '%s' not found", path)
	}

	return
}

// replace custom config vars ($APP_HOME) from the actual values
// and also set default values
func replaceCustomVars(cfg *Configuration) {
	appPath := files.GetAppPath()

	// general defaults
	if cfg.Server.Port == "" {
		cfg.Server.Port = "8080"
	}

	// default DB location and name
	if cfg.Database.Filename == "" {
		cfg.Database.Filename = "$APP_HOME/blog.db"
	}

	// default template settings
	if cfg.Template.Base == "" {
		cfg.Template.Base = "$APP_HOME/templates"
	}
	if cfg.Template.ProcessedOK == "" {
		cfg.Template.ProcessedOK = "$APP_HOME/templates/ok"
	}
	if cfg.Template.ProcessedError == "" {
		cfg.Template.ProcessedError = "$APP_HOME/templates/error"
	}
	if cfg.Template.CheckCycle == 0 {
		cfg.Template.CheckCycle = 30 // 30 seconds
	}

	cfg.Database.Filename = path.Clean(strings.Replace(cfg.Database.Filename, "$APP_HOME", appPath, -1))
	cfg.Template.Base = path.Clean(strings.Replace(cfg.Template.Base, "$APP_HOME", appPath, -1))
	cfg.Template.ProcessedOK = path.Clean(strings.Replace(cfg.Template.ProcessedOK, "$APP_HOME", appPath, -1))
	cfg.Template.ProcessedError = path.Clean(strings.Replace(cfg.Template.ProcessedError, "$APP_HOME", appPath, -1))

}

// load settings from yaml file
func loadSettingsFromFile(path string) (cfg *Configuration, err error) {
	cfg = new(Configuration)

	fileContent, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading configuration file: %s", err)
	}

	if err := yaml.Unmarshal(fileContent, cfg); err != nil {
		return nil, fmt.Errorf("error parsing configuration file: %s", err)
	}

	return
}
