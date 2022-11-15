#! /usr/bin/env bash
function main(){
    echo "build app..."
    go build -ldflags="-X main.version=1.0.0" -o ./build/dockygo ./cmd/...
    cp "./build/dockygo" "$(get_go_bin)"
}

function get_go_bin(){
    gopath="$(go env | grep GOPATH | cut -d "=" -f2 | tr -d '"')"
    echo "$gopath/bin"
}

main "$@"

