package main

import (
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"os"
)

type GOProxy string

const (
	GlobalGoProxy GOProxy = "global"
	ChinaGoProxy  GOProxy = "china"
	GolangProxy   GOProxy = "golang"
)

type GolangOpts struct {
	Version string  `yaml:"version"`
	Proxy   GOProxy `yaml:"proxy"`
}
type Registry struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

type ImageSettings struct {
	Compress   bool   `yaml:"compress"`
	Cache      bool   `yaml:"cache"`
	HTTPProxy  string `yaml:"http_proxy"`
	HTTPSProxy string `yaml:"https_proxy"`
	NOProxy    string `yaml:"no_proxy"`
}
type ImageConfig struct {
	Name        string            `yaml:"name"`
	UserName    string            `yaml:"username"`
	Environment map[string]string `yaml:"environment"`
	Settings    ImageSettings     `yaml:"settings"`
	Golang      GolangOpts        `yaml:"golang"`
}

type DockerImageConfigFile struct {
	Registries    []Registry  `yaml:"registries"`
	ImageSettings ImageConfig `yaml:"image"`
}

func ReadYamlConfigFile(path string) (*DockerImageConfigFile, error) {
	config := &DockerImageConfigFile{}
	yamlFile, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer yamlFile.Close()
	yamlBytes, err := io.ReadAll(yamlFile)
	if err != nil {
		return nil, err
	}
	if err = yaml.Unmarshal(yamlBytes, &config); err != nil {
		return nil, err
	}
	return config, nil
}
