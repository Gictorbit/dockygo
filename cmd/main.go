package main

import (
	"fmt"
	clipkg "gopkg.in/alecthomas/kingpin.v2"
	"log"
	"os"
)

// TODO add version
// TODO add validation http_proxy url,goproxy,registry,go version
// TODO add build date
// TODO add git version
// TODO add go verison

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
		fmt.Printf("%+v\n", dockyGoCmd.BuildCMD)

	case dockyGoCmd.ReleaseCMD.Command.FullCommand():
		fmt.Printf("%+v\n", dockyGoCmd.ReleaseCMD)
	}

}
