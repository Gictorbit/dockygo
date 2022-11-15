package main

import (
	clipkg "gopkg.in/alecthomas/kingpin.v2"
	"log"
	"os/exec"
	"regexp"
)

type BuildCommand struct {
	Command    *clipkg.CmdClause
	ImageName  *string
	Registry   *string
	UserName   *string
	GoProxy    *string
	Compress   *bool
	DockerFile *string
	HTTPProxy  *string
	HTTPSProxy *string
	NOProxy    *string
	Cache      *bool
	GoVersion  *string
	Tags       *[]string
	LatestTag  *bool
}

type ReleaseCommand struct {
	Command   *clipkg.CmdClause
	ImageName *string
	Registry  *string
	UserName  *string
	LatestTag *bool
	Tags      *[]string
}

type DockyGoCmd struct {
	BuildCMD    *BuildCommand
	ReleaseCMD  *ReleaseCommand
	Application *clipkg.Application
}

func MakeCommandLine() *DockyGoCmd {
	app := clipkg.New("dockygo", "docker image builder for go micro services")
	buildCommand := app.Command("build", "build docker image")
	releaseCommand := app.Command("release", "release built image to registry")
	app.HelpFlag.Short('h')

	buildCmd := &BuildCommand{
		Command:    buildCommand,
		ImageName:  buildCommand.Flag("name", "name of image").Short('n').String(),
		Registry:   buildCommand.Flag("registry", "name or url of registry").Short('r').Required().String(),
		UserName:   buildCommand.Flag("user", "registry username").Short('u').String(),
		GoProxy:    buildCommand.Flag("goproxy", "set go mod proxy").Short('g').Default("golang").String(),
		Compress:   buildCommand.Flag("compress", "reduce go binary using upx").Short('c').Default("false").Bool(),
		DockerFile: buildCommand.Flag("dockerfile", "path of docker file").Short('f').Default("./Dockerfile").ExistingFile(),
		HTTPProxy:  buildCommand.Flag("http-proxy", "set http proxy").Envar("http_proxy").String(),
		HTTPSProxy: buildCommand.Flag("https-proxy", "set https proxy").Envar("https_proxy").String(),
		NOProxy:    buildCommand.Flag("no-proxy", "set no proxy").Envar("no_proxy").String(),
		Cache:      buildCommand.Flag("cache", "enable docker cache").Default("false").Bool(),
		GoVersion:  buildCommand.Flag("goversion", "specify go version").Default(GetGoVersion()).String(),
		Tags:       buildCommand.Flag("tags", "docker image tags").Short('t').Strings(),
		LatestTag:  buildCommand.Flag("latest", "build latest tag for image").Short('l').Default("true").Bool(),
	}
	releaseCmd := &ReleaseCommand{
		Command:   releaseCommand,
		Tags:      releaseCommand.Flag("tags", "release image tags").Short('t').Strings(),
		LatestTag: releaseCommand.Flag("latest", "release latest tag for image").Short('l').Default("true").Bool(),
		ImageName: releaseCommand.Flag("name", "name of image").Short('n').String(),
		Registry:  releaseCommand.Flag("registry", "name or url of registry").Short('r').Required().String(),
		UserName:  releaseCommand.Flag("user", "registry username").Short('u').String(),
	}
	return &DockyGoCmd{
		Application: app,
		BuildCMD:    buildCmd,
		ReleaseCMD:  releaseCmd,
	}
}

func GetGoVersion() string {
	versionRegex, err := regexp.Compile("(\\d+\\.)?(\\d+\\.)?(\\*|\\d+)")
	if err != nil {
		log.Fatal(err)
	}
	output, err := exec.Command("go", "list", "-f", "\"{{.GoVersion}}\"", "-m").Output()
	version := versionRegex.FindString(string(output))

	if len(version) > 0 && err == nil {
		return version
	}
	output, err = exec.Command("go", "version", "|", "awk '{print $3}'").Output()
	if err != nil {
		return ""
	}
	goVersion := versionRegex.FindString(string(output))
	if len(goVersion) > 0 && err == nil {
		return goVersion
	}
	return ""
}
