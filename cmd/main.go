package main

import (
	"fmt"
	clipkg "gopkg.in/alecthomas/kingpin.v2"
	"os"
)

var (
	GoProxyServers = map[string]string{
		"golang": "https://proxy.golang.org",
		"china":  "https://goproxy.cn",
		"global": "https://goproxy.io",
	}
)

// TODO add version
// TODO add validation http_proxy url,goproxy,registry,go version
// TODO add build date
// TODO add git version
// TODO add go verison

var (
	dockygo         = clipkg.New("dockygo", "docker image builder for go micro services")
	buildCommand    = dockygo.Command("build", "build docker image")
	buildRegistry   = buildCommand.Flag("registry", "name or url of registry").Short('r').Required().String()
	buildUserName   = buildCommand.Flag("user", "registry username").Short('u').Required().String()
	buildGoProxy    = buildCommand.Flag("goproxy", "set go mod proxy").Short('g').Default("golang").String()
	buildCompress   = buildCommand.Flag("compress", "reduce go binary using upx").Short('c').Bool()
	buildDockerFile = buildCommand.Flag("dockerfile", "path of docker file").Short('f').Default("./Dockerfile").ExistingFile()
	buildHTTPProxy  = buildCommand.Flag("http-proxy", "set http proxy").Envar("http_proxy").String()
	buildHTTPSProxy = buildCommand.Flag("https-proxy", "set https proxy").Envar("https_proxy").String()
	buildNoProxy    = buildCommand.Flag("no-proxy", "set no proxy").Envar("no_proxy").String()
	buildCache      = buildCommand.Flag("cache", "enable docker cache").Bool()
	buildGoVersion  = buildCommand.Flag("goversion", "specify go version").Default("1.19").String()
	buildTag        = buildCommand.Flag("tags", "docker image tags").Short('t').Default("latest").Strings()
	releaseCommand  = dockygo.Command("release", "release built image to registry")
)

func main() {
	dockygo.HelpFlag.Short('h')
	switch clipkg.MustParse(dockygo.Parse(os.Args[1:])) {
	// build docker image
	case buildCommand.FullCommand():
		serverAddr, ok := GoProxyServers[*buildGoProxy]
		if !ok {
			clipkg.Errorf("invalid go proxy")
			os.Exit(1)
		}
		fmt.Println("build")
		fmt.Println("registry:", *buildRegistry)
		fmt.Println("user:", *buildUserName)
		fmt.Println("goproxy:", *buildGoProxy, serverAddr)
		fmt.Println("compress:", *buildCompress)
		fmt.Println("dockerfile:", *buildDockerFile)
		fmt.Println("http proxy:", *buildHTTPProxy)
		fmt.Println("https proxy:", *buildHTTPSProxy)
		fmt.Println("no proxy:", *buildNoProxy)
		fmt.Println("cache", *buildCache)
		fmt.Println("go version:", *buildGoVersion)
		fmt.Println("tags", *buildTag)

	case releaseCommand.FullCommand():
		fmt.Println("release")
	}
	config, err := ReadYamlConfigFile("./Dockerimg.yaml")
	if err != nil {
		panic(err)
	}
	fmt.Println(config)
}
