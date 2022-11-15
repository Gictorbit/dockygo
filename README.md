# Dockygo
docker image builder for go micro services
## Installation
you can build and install `dockygo` on your linux machine using builder script. <br/>
1- first clone the repository
```shell
git clone "https://github.com/Gictorbit/dockygo.git"
```
2- execute builder script
```shell
cd dockygo && chmod +x build.sh && ./build.sh
```
## Config File Example
by default `dockygo` looks for configuration file in current directory like `Dockerimg.yaml`. <br/>
users can overwrite file options using command line flags
```yaml
image:
  name: dockygo-api
  username: dockygo
  environment:
    COMPANY_HOST: "github.com/dockygo"
    GITHUB_TOKEN: "${GITHUB_TOKEN}" #looks for env variable
  settings:
    compress: true
    latest: false
    cache: false
    http_proxy:
    https_proxy:
    no_proxy:
  golang:
    proxy: global #[golang,global,china] default:global
    version: 1.19 #default: current installed version
registries:
  - name: local
    url: reg.dockygo.com
```
## Usage
```shell
usage: dockygo [<flags>] <command> [<args> ...]

docker image builder for go micro services

Flags:
  -h, --help  Show context-sensitive help (also try --help-long and --help-man).

Commands:
  help [<command>...]
    Show help.

  build --registry=REGISTRY [<flags>]
    build docker image

  release --registry=REGISTRY [<flags>]
    release built image to registry
```

### Build
use build command to build your docker image 
```shell
build docker image

Flags:
  -h, --help                     Show context-sensitive help (also try --help-long and --help-man).
  -n, --name=""                  name of image
  -r, --registry=REGISTRY        name or url of registry
  -u, --user=USER                registry username
  -g, --goproxy="golang"         set go mod proxy
  -c, --compress                 reduce go binary using upx
  -f, --dockerfile=./Dockerfile  path of docker file
      --http-proxy=HTTP-PROXY    set http proxy
      --https-proxy=HTTPS-PROXY  set https proxy
      --no-proxy=NO-PROXY        set no proxy
      --cache                    enable docker cache
      --goversion="1.19"         specify go version
  -t, --tag="v1.0.1"             docker image tag
  -l, --latest                   build latest tag for image
```
### Release
use release command to push your built docker image to registry
```shell
usage: dockygo release --registry=REGISTRY [<flags>]

release built image to registry

Flags:
  -h, --help               Show context-sensitive help (also try --help-long and --help-man).
  -t, --tag="v1.0.1"       release image tag
  -l, --latest             release latest tag for image
  -n, --name=NAME          name of image
  -r, --registry=REGISTRY  name or url of registry
  -u, --user=USER          registry username

```