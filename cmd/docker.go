package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

type DockerBuildOptions struct {
	Tags       []string
	RemoteAddr string
	BuildArg   map[string]string
	BuildX     bool
	Dockerfile string
	Cache      bool
}
type DockerReleaseOptions struct {
	Tags       []string
	RemoteAddr string
}

func (dbo *DockerBuildOptions) BuildCommand() *exec.Cmd {
	flags := make([]string, 0)
	if dbo.BuildX {
		flags = append(flags, "buildx", "build")
	} else {
		flags = append(flags, "build")
	}
	for _, tg := range dbo.Tags {
		tag := fmt.Sprintf("%s:%s", dbo.RemoteAddr, tg)
		flags = append(flags, "-t", tag)
	}
	for arg, value := range dbo.BuildArg {
		param := fmt.Sprintf("%s=%s", arg, value)
		flags = append(flags, "--build-arg", param)
	}
	if !dbo.Cache {
		flags = append(flags, "--no-cache")
	}
	flags = append(flags, "--network", "host", ".")
	return exec.Command("docker", flags...)
}

func BuildDockerImage(opts DockerBuildOptions) error {
	cmd := opts.BuildCommand()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		return err
	}

	return nil
}

func GetFullRemoteAddr(config *DockerImageConfigFile) string {
	return fmt.Sprintf("%s/%s/%s", config.RemoteAddr, config.ImageSettings.UserName, config.ImageSettings.Name)
}

func PrintBuildConfig(config *DockerImageConfigFile) {
	infos := map[string]any{
		"RemoteAddr": config.RemoteAddr,
		"Username":   config.ImageSettings.UserName,
		"ImageName":  config.ImageSettings.Name,
		"Compress":   config.ImageSettings.Settings.Compress,
		"Cache":      config.ImageSettings.Settings.Cache,
		"GoVersion":  config.ImageSettings.Golang.Version,
		"Tags":       config.Tags,
		"HTTPProxy":  config.ImageSettings.Settings.HTTPProxy,
		"HTTPSProxy": config.ImageSettings.Settings.HTTPSProxy,
		"NOProxy":    config.ImageSettings.Settings.NOProxy,
	}
	for key, value := range infos {
		log.Printf("%v:\t%v\n", key, value)
	}
	log.Println("environments:")
	for key, value := range config.ImageSettings.Environment {
		fmt.Printf("\t%s: %s\n", key, value)
	}
	log.Println("images:")
	for _, tag := range config.Tags {
		fmt.Printf("\t%s:%s\n", GetFullRemoteAddr(config), tag)
	}
}

func PrintReleaseConfig(config *DockerImageConfigFile) {
	infos := map[string]any{
		"RemoteAddr": config.RemoteAddr,
		"Username":   config.ImageSettings.UserName,
		"ImageName":  config.ImageSettings.Name,
		"Tags":       config.Tags,
	}
	for key, value := range infos {
		log.Printf("%v:\t%v\n", key, value)
	}
	log.Println("images:")
	for _, tag := range config.Tags {
		fmt.Printf("\t%s:%s\n", GetFullRemoteAddr(config), tag)
	}
}

func PushDockerImage(opts DockerReleaseOptions) error {
	for _, tag := range opts.Tags {
		tg := fmt.Sprintf("%s:%s", opts.RemoteAddr, tag)
		cmd := exec.Command("docker", "push", tg)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Start(); err != nil {
			return err
		}
		if err := cmd.Wait(); err != nil {
			return err
		}
	}
	return nil
}
