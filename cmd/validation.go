package main

import (
	"errors"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var (
	ErrInvalidUserName  = errors.New("empty or invalid username")
	ErrInvalidImageName = errors.New("empty or invalid image name")
	ErrInvalidGoProxy   = errors.New("invalid go proxy")
	ErrInvalidGoVersion = errors.New("empty or invalid go version")
	ErrEmptyTags        = errors.New("no tags specified")
)

func ValidateBuild(config *DockerImageConfigFile, cmd *BuildCommand) error {
	//validate imageName
	if len(cmd.ImageName) > 0 {
		config.ImageSettings.Name = cmd.ImageName
	}
	if len(config.ImageSettings.Name) == 0 {
		return ErrInvalidImageName
	}
	//validate username
	if len(cmd.UserName) > 0 {
		config.ImageSettings.UserName = cmd.UserName
	}
	if len(config.ImageSettings.UserName) == 0 {
		return ErrInvalidUserName
	}
	for _, reg := range config.Registries {
		if reg.Name != cmd.Registry {
			continue
		}

	}
	//validate go proxy
	if len(config.ImageSettings.Golang.Proxy) == 0 {
		config.ImageSettings.Golang.Proxy = GolangProxy
	}
	_, configOk := GoProxyServers[config.ImageSettings.Golang.Proxy]
	_, inputOk := GoProxyServers[GOProxy(cmd.GoProxy)]
	if !inputOk || !configOk {
		return ErrInvalidGoProxy
	}
	if GOProxy(cmd.GoProxy) != GolangProxy {
		config.ImageSettings.Golang.Proxy = GOProxy(cmd.GoProxy)
	}
	//validate registry
	for _, reg := range config.Registries {
		if reg.Name == cmd.Registry {
			config.RemoteAddr = reg.URL
		}
	}
	if len(config.RemoteAddr) == 0 {
		config.RemoteAddr = cmd.Registry
	}
	//validate compress
	if cmd.Compress == true {
		config.ImageSettings.Settings.Compress = true
	}
	//validate http proxy
	if len(cmd.HTTPProxy) > 0 {
		config.ImageSettings.Settings.HTTPProxy = cmd.HTTPProxy
	}
	if len(cmd.HTTPSProxy) > 0 {
		config.ImageSettings.Settings.HTTPSProxy = cmd.HTTPSProxy
	}
	if len(cmd.NOProxy) > 0 {
		config.ImageSettings.Settings.NOProxy = cmd.NOProxy
	}
	if cmd.Cache == true {
		config.ImageSettings.Settings.Cache = true
	}
	//validate latest
	if cmd.LatestTag {
		config.ImageSettings.Settings.Latest = true
	}
	//validate go version
	if len(config.ImageSettings.Golang.Version) == 0 && len(cmd.GoVersion) > 0 ||
		len(cmd.GoVersion) > 0 && cmd.GoVersion != GetGoVersion() {
		config.ImageSettings.Golang.Version = cmd.GoVersion
	}
	if len(config.ImageSettings.Golang.Version) == 0 {
		return ErrInvalidGoVersion
	}
	// validate tags
	if len(cmd.Tag) > 0 {
		config.Tags = append(config.Tags, cmd.Tag)
	}
	if gitVersion := GetRepoTagVersion(); len(gitVersion) > 0 && len(cmd.Tag) == 0 {
		config.Tags = append(config.Tags, gitVersion)
	}
	if config.ImageSettings.Settings.Latest {
		config.Tags = append(config.Tags, "latest")
	}
	if len(config.Tags) == 0 {
		return ErrEmptyTags
	}
	return nil
}

func AddBuildArgs(config *DockerImageConfigFile) {
	config.ImageSettings.Environment["BUILD_DATE"] = time.Now().UTC().String()
	config.ImageSettings.Environment["GO_VERSION"] = config.ImageSettings.Golang.Version
	config.ImageSettings.Environment["COMPRESS"] = strconv.FormatBool(config.ImageSettings.Settings.Compress)
	config.ImageSettings.Environment["USERNAME"] = config.ImageSettings.UserName
	config.ImageSettings.Environment["IMAGE_NAME"] = config.ImageSettings.Name
	config.ImageSettings.Environment["GOPROXYURL"] = GoProxyServers[config.ImageSettings.Golang.Proxy]
	if len(config.ImageSettings.Settings.NOProxy) > 0 {
		config.ImageSettings.Environment["NO_PROXY"] = config.ImageSettings.Settings.NOProxy
	}
	if len(config.ImageSettings.Settings.HTTPProxy) > 0 {
		config.ImageSettings.Environment["HTTP_PROXY"] = config.ImageSettings.Settings.HTTPProxy
	}
	if len(config.ImageSettings.Settings.HTTPSProxy) > 0 {
		config.ImageSettings.Environment["HTTPS_PROXY"] = config.ImageSettings.Settings.HTTPSProxy
	}
	config.ImageSettings.Environment["VERSION"] = GetRepoTagVersion()
}

func GetRepoTagVersion() string {
	args1 := []string{"describe", "--tags", "--exact-match", "2>/dev/null"}
	args2 := []string{"symbolic-ref", "-q", "--short", "HEAD"}
	output, err := exec.Command("git", args1...).Output()
	if err == nil && len(string(output)) > 0 {
		return strings.TrimSuffix(string(output), "\n")
	}
	output, err = exec.Command("git", args2...).Output()
	if err == nil && len(string(output)) > 0 {
		return strings.TrimSuffix(string(output), "\n")
	}
	return ""
}
