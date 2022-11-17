package main

import (
	clipkg "gopkg.in/alecthomas/kingpin.v2"
	"log"
	"os/exec"
	"regexp"
)

var Version = "development"

type BuildCommand struct {
	Command    *clipkg.CmdClause
	ImageName  string
	Registry   string
	UserName   string
	GoProxy    string
	Compress   bool
	DockerFile string
	HTTPProxy  string
	HTTPSProxy string
	NOProxy    string
	Cache      bool
	GoVersion  string
	Tag        string
	LatestTag  bool
}

type ReleaseCommand struct {
	Command   *clipkg.CmdClause
	ImageName string
	Registry  string
	UserName  string
	LatestTag bool
	Tag       string
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
	app.Version(Version)
	app.VersionFlag.Short('v')
	gitTag := GetRepoTagVersion()
	buildCmd := &BuildCommand{Command: buildCommand}
	buildCommand.Flag("name", "name of image").Short('n').Default("").StringVar(&buildCmd.ImageName)
	buildCommand.Flag("registry", "name or url of registry").Short('r').Required().StringVar(&buildCmd.Registry)
	buildCommand.Flag("user", "registry username").Short('u').StringVar(&buildCmd.UserName)
	buildCommand.Flag("goproxy", "set go mod proxy").Short('g').Default("golang").StringVar(&buildCmd.GoProxy)
	buildCommand.Flag("compress", "reduce go binary using upx").Short('c').Default("false").BoolVar(&buildCmd.Compress)
	buildCommand.Flag("dockerfile", "path of docker file").Short('f').Default("./Dockerfile").ExistingFileVar(&buildCmd.DockerFile)
	buildCommand.Flag("http-proxy", "set http proxy").Envar("http_proxy").StringVar(&buildCmd.HTTPProxy)
	buildCommand.Flag("https-proxy", "set https proxy").Envar("https_proxy").StringVar(&buildCmd.HTTPSProxy)
	buildCommand.Flag("no-proxy", "set no proxy").Envar("no_proxy").StringVar(&buildCmd.NOProxy)
	buildCommand.Flag("cache", "enable docker cache").Default("false").BoolVar(&buildCmd.Cache)
	buildCommand.Flag("goversion", "specify go version").Default(GetGoVersion()).StringVar(&buildCmd.GoVersion)
	buildCommand.Flag("tag", "docker image tag").Short('t').Default(gitTag).StringVar(&buildCmd.Tag)
	buildCommand.Flag("latest", "build latest tag for image").Short('l').Default("false").BoolVar(&buildCmd.LatestTag)

	releaseCmd := &ReleaseCommand{Command: releaseCommand}
	releaseCommand.Flag("tag", "release image tag").Short('t').Default(gitTag).StringVar(&releaseCmd.Tag)
	releaseCommand.Flag("latest", "release latest tag for image").Short('l').Default("false").BoolVar(&releaseCmd.LatestTag)
	releaseCommand.Flag("name", "name of image").Short('n').StringVar(&releaseCmd.ImageName)
	releaseCommand.Flag("registry", "name or url of registry").Short('r').Required().StringVar(&releaseCmd.Registry)
	releaseCommand.Flag("user", "registry username").Short('u').StringVar(&releaseCmd.UserName)

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
