package main

import (
	"fmt"
	clipkg "gopkg.in/alecthomas/kingpin.v2"
	"log"
	"os"
)

func main() {
	dockyGoCmd := MakeCommandLine()
	parsedCmd := clipkg.MustParse(dockyGoCmd.Application.Parse(os.Args[1:]))
	path := GetYamlConfigFilePath()
	yamlConfig := &DockerImageConfigFile{}
	if path != "" {
		var err error
		yamlConfig, err = ReadYamlConfigFile(path)
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Printf("%+v\n*********\n", yamlConfig)
	switch parsedCmd {
	// build docker image
	case dockyGoCmd.BuildCMD.Command.FullCommand():
		if err := ValidateBuild(yamlConfig, dockyGoCmd.BuildCMD); err != nil {
			log.Fatal(err)
		}
		AddBuildArgs(yamlConfig)
		dockerOpts := DockerBuildOptions{
			BuildX:     true,
			Tags:       yamlConfig.Tags,
			RemoteAddr: GetFullTag(yamlConfig),
			Dockerfile: *dockyGoCmd.BuildCMD.DockerFile,
			Cache:      yamlConfig.ImageSettings.Settings.Cache,
			BuildArg:   yamlConfig.ImageSettings.Environment,
		}
		if err := BuildDockerImage(dockerOpts); err != nil {
			log.Fatal(err)
		}

	case dockyGoCmd.ReleaseCMD.Command.FullCommand():
		fmt.Printf("%+v\n", dockyGoCmd.ReleaseCMD)
	}

}
