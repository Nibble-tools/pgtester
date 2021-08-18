package internal

import (
	"flag"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"mannemsolutions/pgtester/pkg/pg"
	"os"
	"path/filepath"
)

/*
 * This module reads the config file and returns a config object with all entries from the config yaml file.
 */

const (
	envConfName     = "PGTESTS"
	defaultConfFile = "./tests.yaml"
)

type Tests []Test

type Test struct {
	Query  string            `yaml:"query"`
	Results pg.OneFieldResults `yaml:"results"`
}

type Config struct {
	Debug bool			`yaml:"debug"`
	Tests Tests         `yaml:"tests"`
	DSN   pg.Dsn        `yaml:"dsn"`
}

func NewConfig() (config Config, err error) {
	var configFile string
	var debug bool
	flag.StringVar(&configFile, "f", os.Getenv(envConfName), "Specify file with tests")
	flag.BoolVar(&debug, "d", false, "Add debugging output")
	flag.Parse()
	 if configFile == "" {
		 configFile = defaultConfFile
	 }
	configFile, err = filepath.EvalSymlinks(configFile)
	if err != nil {
		config.Debug = config.Debug || debug
		return config, err
	}

	// This only parsed as yaml, nothing else
	// #nosec
	yamlConfig, err := ioutil.ReadFile(configFile)
	if err != nil {
		return config, err
	}
	err = yaml.Unmarshal(yamlConfig, &config)
	return config, err
}
