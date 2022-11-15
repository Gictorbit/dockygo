package main

import (
	"errors"
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type GOProxy string

const (
	GlobalGoProxy GOProxy = "global"
	ChinaGoProxy  GOProxy = "china"
	GolangProxy   GOProxy = "golang"
)

var (
	GoProxyServers = map[GOProxy]string{
		GolangProxy:   "https://proxy.golang.org",
		ChinaGoProxy:  "https://goproxy.cn",
		GlobalGoProxy: "https://goproxy.io",
	}
	ConfigFiles = []string{"Dockerimg.yaml", "Dockerimg.yml"}
)

type GolangOpts struct {
	Version string  `yaml:"version,omitempty"`
	Proxy   GOProxy `yaml:"proxy,omitempty"`
}
type Registry struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

type ImageSettings struct {
	Compress   bool   `yaml:"compress,omitempty"`
	Cache      bool   `yaml:"cache,omitempty"`
	HTTPProxy  string `yaml:"http_proxy,omitempty"`
	HTTPSProxy string `yaml:"https_proxy,omitempty"`
	NOProxy    string `yaml:"no_proxy,omitempty"`
	Latest     bool   `yaml:"latest,omitempty"`
}
type ImageConfig struct {
	Name        string            `yaml:"name,omitempty"`
	UserName    string            `yaml:"username,omitempty"`
	Environment map[string]string `yaml:"environment,omitempty"`
	Settings    ImageSettings     `yaml:"settings,omitempty"`
	Golang      GolangOpts        `yaml:"golang,omitempty"`
}

type DockerImageConfigFile struct {
	RemoteAddr    string
	Tags          []string
	Registries    []Registry  `yaml:"registries,omitempty"`
	ImageSettings ImageConfig `yaml:"image,omitempty"`
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
	r, err := regexp.Compile("\\${([^}]+)}")
	if err != nil {
		return nil, err
	}
	replacer := strings.NewReplacer("$", "", "{", "", "}", "")
	for key, value := range config.ImageSettings.Environment {
		if r.MatchString(value) {
			envVar := replacer.Replace(r.FindString(value))
			config.ImageSettings.Environment[key] = os.Getenv(envVar)
		}
	}
	return config, nil
}

func GetYamlConfigFilePath() string {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	for _, cfgFile := range ConfigFiles {
		path := filepath.Join(pwd, cfgFile)
		if _, existsErr := os.Stat(path); !errors.Is(existsErr, os.ErrNotExist) {
			return path
		}
	}
	return ""
}
