package main

import (
	"bufio"
	"fmt"
	"log"
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
	stderr, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		return err
	}
	scanner := bufio.NewScanner(stderr)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
	}
	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
}

func GetFullTag(config *DockerImageConfigFile) string {
	return fmt.Sprintf("%s/%s/%s", config.RemoteAddr, config.ImageSettings.UserName, config.ImageSettings.Name)
}

func PrintBuildConfig(config *DockerImageConfigFile) {
	log.Printf("remoteAddr: %s\n", config.RemoteAddr)
	log.Printf("image name: %s\n", config.ImageSettings.Name)
	log.Printf("compress: %v\n", config.ImageSettings.Settings.Compress)
	log.Printf("cache: %v\n", config.ImageSettings.Settings.Cache)
	log.Printf("go version: %s\n", config.ImageSettings.Golang.Version)
	log.Printf("tags: %v\n", config.Tags)
	log.Println("environments:")
	for key, value := range config.ImageSettings.Environment {
		fmt.Printf("\t%s: %s\n", key, value)
	}
}
