#!/usr/bin/env bash

function main() {
    initialize
    read_config_file
    parseargs "$@"
    read_env
    validate_args

    message "service name = $SERVICE_NAME"
    message "build date = $BUILD_DATE"
    message "remote address = $REMOTE_ADDR"
    message "compress = $COMPRESS"
    message "go version = $GOLANG_VERSION"
    message "version = $GIT_VERSION"
    message "company host = $COMPANY_HOST"
    message "company name = $COMPANY_NAME"
    message "go proxy = $GOPROXY $GOPROXYURL"
    message "http proxy= $HTTP_PROXY"
    message "https proxy= $HTTPS_PROXY"
    message "no proxy= $NO_PROXY"
    message "cache = $CACHE"

    build_service_image
    [ "$BASECOMMAND" == "release" ] && push_docker_image
}

function read_config_file(){

    CONFIG_FILE="Dockerimg.yaml"
    if [ -f "Dockerimg.yaml" ]; then
        CONFIG_FILE="Dockerimg.yaml"
    else
        return 0
    fi
    config_content="$(cat $CONFIG_FILE)"

    SERVICE_NAME="$(echo "$config_content" | yq -r '.service.name' || echo "")"
    [ -z "$SERVICE_NAME" ] || [ "$SERVICE_NAME" == "null" ] && SERVICE_NAME=""

    COMPANY_HOST="$(echo "$config_content" | yq -r '.service.company_host' || echo "")"
    [ -z "$COMPANY_HOST" ] || [ "$COMPANY_HOST" == "null" ] && COMPANY_HOST=""

    COMPANY_NAME="$(echo "$config_content" | yq -r '.service.company_name' || echo "")"
    [ -z "$COMPANY_NAME" ] || [ "$COMPANY_NAME" == "null" ] && COMPANY_NAME=""

    GOPROXY="$(echo "$config_content" | yq -r '.service.settings.go_proxy' || echo "global")"
    [ -z "$GOPROXY" ] || [ "$GOPROXY" == "null" ] && GOPROXY="global"

    COMPRESS="$(echo "$config_content" | yq -r '.service.settings.compress' || echo "false")"
    [ -z "$COMPRESS" ] || [ "$COMPRESS" == "null" ] && COMPRESS="false"

    CACHE="$(echo "$config_content" | yq -r '.service.settings.cache' || echo "false")"
    [ -z "$CACHE" ] || [ "$CACHE" == "null" ] && CACHE="false"

    GOLANG_VERSION="$(echo "$config_content" | yq -r '.service.settings.go_version' || echo "")"
    [ -z "$GOLANG_VERSION" ] || [ "$GOLANG_VERSION" == "null" ] && GOLANG_VERSION="$(go_module_version)"

    HTTP_PROXY="$(echo "$config_content" | yq -r '.service.settings.http_proxy' || echo "")"
    [ -z "$HTTP_PROXY" ] || [ "$HTTP_PROXY" == "null" ] && HTTP_PROXY=""

    HTTPS_PROXY="$(echo "$config_content" | yq -r '.service.settings.https_proxy' || echo "")"
    [ -z "$HTTPS_PROXY" ] || [ "$HTTPS_PROXY" == "null" ] && HTTPS_PROXY=""

    NO_PROXY="$(echo "$config_content" | yq -r '.service.settings.no_proxy' || echo "")"
    [ -z "$NO_PROXY" ] || [ "$NO_PROXY" == "null" ] && NO_PROXY=""

    defined_regs="$(echo "$config_content" | yq -r '.service.registries[]| "[\"" + (.name) + "\"]=" + (.url)' || echo "")"
    [ -z "$defined_regs" ] || [ "$defined_regs" == "null" ] && defined_regs=""
    declare -gA registries="($defined_regs)"
}

function initialize() {
    DATE_FMT="+%FT%T%z"
    GOLANG_VERSION="$(go_module_version)"
    COMPANY_HOST=""
    COMPANY_NAME=""
    SERVICE_NAME=""
    REMOTE_ADDR=""
    COMPRESS="false"
    GOPROXY="global"
    HTTP_PROXY=""
    HTTPS_PROXY=""
    NO_PROXY=""
    GIT_VERSION=$(git describe --tags --exact-match 2>/dev/null || git symbolic-ref -q --short HEAD)
    BUILD_DATE=$(date "$DATE_FMT")
    CACHE="false"

    declare -gA basecmd_list
    basecmd_list["help"]="prints help"
    basecmd_list["build"]="build docker image"
    basecmd_list["release"]="build and release docker image to registry"

    declare -gA flags_list
    flags_list["-r,--registry"]="registry name or url"
    flags_list["-c,--compress"]="reduce docker image size using upx"
    flags_list["-g,--goproxy"]="set go module proxy default: global"
    flags_list["-s,--service"]="set service name"
    flags_list["--cname"]="set company name"
    flags_list["--chost"]="set company host"
    flags_list["--http-proxy"]="set http proxy"
    flags_list["--https-proxy"]="set https proxy"
    flags_list["--no-proxy"]="set no proxy"
    flags_list["--cache"]="enable cache for docker build"


    declare -gA go_proxies
    go_proxies["golang"]="https://proxy.golang.org"
    go_proxies["global"]="https://goproxy.io"
    go_proxies["china"]="https://goproxy.cn"
}

function read_env() {
    GITHUB_TOKEN=$(printenv | grep "GITHUB_TOKEN" | cut -d "=" -f 2)
    [ -z "$HTTP_PROXY" ] && HTTP_PROXY="$(printenv | grep http_proxy | cut -d "=" -f2)"
    [ -z "$HTTPS_PROXY" ] && HTTPS_PROXY="$(printenv | grep https_proxy | cut -d "=" -f2)"
    [ -z "$NO_PROXY" ] && NO_PROXY="$(printenv | grep no_proxy | cut -d "=" -f2)"
}

function parseargs() {
    BASECOMMAND=$1
    [[ "$BASECOMMAND" =~ ^((-h)|(--help)|(help))$ ]] && print_help
    shift
    while [ $# -gt 0 ]; do
        case $1 in
        -r | --registry)
            REGISTRY="$2"
            shift
            ;;
        -c | --compress)
            COMPRESS="true"
            ;;
        -g | --goproxy)
            GOPROXY="$2"
            shift
            ;;
        --cname)
            COMPANY_NAME="$2"
            shift
            ;;
        --chost)
            COMPANY_HOST="$2"
            shift
            ;;
        --goversion)
            GOLANG_VERSION="$2"
            shift
            ;;
        -s | --service)
            SERVICE_NAME="$2"
            shift
            ;;
        --hproxy)
            HTTP_PROXY="$2"
            shift
            ;;
        --hsproxy)
            HTTPS_PROXY="$2"
            shift
            ;;
        --noproxy)
            NO_PROXY="$2"
            shift
            ;;
        --cache)
            CACHE="true"
            ;;
        -h | --help) print_help ;;
        *) print_help ;;
        esac
        shift
    done
}

function validate_args() {
    #validate base command
    local found_base_cmd=0
    for cmd in "${!basecmd_list[@]}"; do
        if [ "$cmd" == "$BASECOMMAND" ]; then
            found_base_cmd=1
            break
        fi
    done
    if [ $found_base_cmd -eq 0 ]; then
        error "Invalid base command: $BASECOMMAND"
    fi
    #validate registry
    REMOTE_ADDR="$REGISTRY"
    for regkey in "${!registries[@]}"; do
        if [ "$regkey" == "$REGISTRY" ]; then
            REMOTE_ADDR="${registries[$regkey]}"
            break
        fi
    done
    [ -z "$REMOTE_ADDR" ] && error "Invalid registry: $REGISTRY"

    #validate go proxy
    if [ -n "$GOPROXY" ]; then
        local gproxy_found=0
        for gpkey in "${!go_proxies[@]}"; do
            if [ "$gpkey" == "$GOPROXY" ]; then
                gproxy_found=1
                break
            fi
        done
        if [ $gproxy_found -eq 0 ]; then
            error "Invalid go proxy: $GOPROXY"
        fi
        GOPROXYURL=${go_proxies[$GOPROXY]}
    fi
    [ -z "$GITHUB_TOKEN" ] && error "GITHUB_TOKEN env not found"
    [ -z "$COMPANY_NAME" ] && error "company name not found"
    [ -z "$COMPANY_HOST" ] && error "company host not found"
    [ -z "$SERVICE_NAME" ] && error "service name not found"
    [ -z "$GOPROXY" ] && error "go module proxy not found"
    [ "$(echo "$GOLANG_VERSION" | grep -Po '(\d+\.)+\d+')" == "" ] && error "invalid go version"
}

function print_help() {

    local first_reg
    first_reg=$(for reg in "${!registries[@]}"; do echo $reg; done | head -n 1)

    printf "\n%s\n\t%s\n" "Usage:" "$0 [build|release] [-r|--registry <registry_name>] [-c|--compress] [--cache]"
    printf "\n%s\n\t%s\n" "Example:" " $0 build --registry $first_reg --compress --cache"
    printf "\t%s\n" " $0 release --registry $first_reg --compress"
    echo

    print_base_commands
    print_arguments
    print_go_proxies
    exit
}

function print_base_commands() {
    echo "base commands:"
    for cmd in "${!basecmd_list[@]}"; do
        printf "${Yellow}\t%-20s${Reset} ${Cyan}%-20s${Reset}\n" $cmd "${basecmd_list[$cmd]}"
    done
    echo
}

function print_arguments() {
    echo "supported flags:"
    for flg in "${!flags_list[@]}"; do
        printf "${Yellow}\t%-20s${Reset} ${Cyan}%-20s${Reset}\n" $flg "${flags_list[$flg]}"
    done
    echo
}

function print_go_proxies() {
    echo "go proxy list:"
    for proxy in "${!go_proxies[@]}"; do
        printf "${Yellow}\t%-20s${Reset} ${Cyan}%-20s${Reset}\n" $proxy "${go_proxies[$proxy]}"
    done
    echo
}

function go_module_version() {
    local module_version
    module_version="$(go list -f "{{.GoVersion}}" -m || echo "")"
    if [ -z "$module_version" ]; then
        module_version="$(go version | awk '{print $3}')"
    fi
    echo "$module_version"
}

function build_service_image() {
    message "build service docker image..."
    local docker_file

    docker_file="$(pwd)/Dockerfile"
    [ ! -f "$docker_file" ] && echo "docker file not found: $docker_file" && exit 1

    if [ "$CACHE" = "true" ];then
        docker buildx build -f "$docker_file" \
            -t "${REMOTE_ADDR}/${COMPANY_NAME}/${SERVICE_NAME}:${GIT_VERSION}" \
            -t "${REMOTE_ADDR}/${COMPANY_NAME}/${SERVICE_NAME}:latest" \
            --build-arg GITHUB_TOKEN="${GITHUB_TOKEN}" \
            --build-arg BUILD_DATE="${BUILD_DATE}" \
            --build-arg GO_VERSION="${GOLANG_VERSION}" \
            --build-arg COMPRESS="${COMPRESS}" \
            --build-arg COMPANY_HOST="${COMPANY_HOST}" \
            --build-arg GOPROXYURL="${GOPROXYURL}" \
            --build-arg HTTP_PROXY="${HTTP_PROXY}" \
            --build-arg HTTPS_PROXY="${HTTPS_PROXY}" \
            --build-arg NO_PROXY="${HTTPS_PROXY}" \
            --build-arg SERVICE_NAME="${SERVICE_NAME}" \
            --build-arg VERSION="${GIT_VERSION}" --network host .
    else
        docker buildx build -f "$docker_file" \
            -t "${REMOTE_ADDR}/${COMPANY_NAME}/${SERVICE_NAME}:${GIT_VERSION}" \
            -t "${REMOTE_ADDR}/${COMPANY_NAME}/${SERVICE_NAME}:latest" \
            --build-arg GITHUB_TOKEN="${GITHUB_TOKEN}" \
            --build-arg BUILD_DATE="${BUILD_DATE}" \
            --build-arg GO_VERSION="${GOLANG_VERSION}" \
            --build-arg COMPRESS="${COMPRESS}" \
            --build-arg COMPANY_HOST="${COMPANY_HOST}" \
            --build-arg GOPROXYURL="${GOPROXYURL}" \
            --build-arg HTTP_PROXY="${HTTP_PROXY}" \
            --build-arg HTTPS_PROXY="${HTTPS_PROXY}" \
            --build-arg NO_PROXY="${HTTPS_PROXY}" \
            --build-arg SERVICE_NAME="${SERVICE_NAME}" \
            --build-arg VERSION="${GIT_VERSION}" --no-cache --network host .
    fi

}

function push_docker_image() {
    message "release image to $REMOTE_ADDR"
    message "push service docker image..."
    docker push "${REMOTE_ADDR}/${COMPANY_NAME}/${SERVICE_NAME}:${GIT_VERSION}"
    docker push "${REMOTE_ADDR}/${COMPANY_NAME}/${SERVICE_NAME}:latest"
}

function error() {
    echo -e "${Red}Error:${Yellow} $1 ${Reset}"
    print_help
    exit 1
}
function message() {
    echo -e "Info: ${Cyan}$1 ${Reset}"
}
function init_colors() {
    # Reset
    Reset='\033[0m' # Text Reset

    # Regular Colors
    Red='\033[0;31m'    # Red
    Green='\033[0;32m'  # Green
    Yellow='\033[0;33m' # Yellow
    Blue='\033[0;34m'   # Blue
    Purple='\033[0;35m' # Purple
    Cyan='\033[0;36m'   # Cyan

    # Bold
    BRed='\033[1;31m'    # Red
    BGreen='\033[1;32m'  # Green
    BYellow='\033[1;33m' # Yellow
    BBlue='\033[1;34m'   # Blue
    BPurple='\033[1;35m' # Purple
    BCyan='\033[1;36m'   # Cyan
}
init_colors
main "$@"