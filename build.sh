#! /usr/bin/env bash
function main(){
    echo "build app..."
    VERSION="$(git describe --tags --exact-match 2>/dev/null || git symbolic-ref -q --short HEAD)"
    BUILD_DATE="$(date -u +%Y-%m-%dT%H:%M:%S%Z)"

    go build -ldflags="-X main.Version=$VERSION -X main.buildDate=$BUILD_DATE" -o ./build/dockygo ./cmd/...
    cp "./build/dockygo" "$(get_go_bin)"
}

function get_go_bin(){
    gopath="$(go env | grep GOPATH | cut -d "=" -f2 | tr -d '"')"
    echo "$gopath/bin"
}

main "$@"

