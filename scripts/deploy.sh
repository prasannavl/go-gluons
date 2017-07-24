#!/usr/bin/env bash
## Author: Prasanna V. Loganathar

set -e
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )/../" && pwd )"
cd ${dir}

binary_name="nextfirst-core"
dest_path="~/workspace/bin/"

# Build
go get -u github.com/golang/dep/cmd/dep
dep ensure && go build -o ./bin/nextfirst-core

# Deploy
mkdir -p "$dest_path"

wait_for_process_briefly(){
    local pid=$(pidof "$@")
    local i=0
    while kill -0 "$pid"; do
        sleep 0.1s
        i=$((i+1))
        if [ $i -gt 900 ]; then kill -9 "$pid"; fi;
    done
}

wait_for_process_briefly "$binary_name"
mv -f "./bin/${binary_name}" "$dest_path"
nohup "${dest_path}/${binary_name}"