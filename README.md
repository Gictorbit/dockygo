# dockygo
docker image builder for go micro services
## Installation
you can build and install `dockygo` on your linux machine using builder script
1- first clone the repository
```shell
git clone "https://github.com/Gictorbit/dockygo.git"
```
2- execute builder script
```shell
cd dockygo && chmod +x builder.sh && ./build.sh
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