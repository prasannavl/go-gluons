#!/usr/bin/env bash
## Author: Prasanna V. Loganathar

binary_name="nextfirst-core"
deploy_dir="${HOME}/workspace/bin/"
logs_dir="${HOME}/workspace/logs/"

main() {
    set -e
    local dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )/../" && pwd )"    
    cd ${dir}
    build
    deploy
    cd ${deploy_dir}    
    start
}

build() {
    echo "> build: start"
    go get -v -u github.com/golang/dep/cmd/dep || true
    dep ensure || true
    go build -o "./bin/${binary_name}"
    echo "> build: success"
}

deploy() {
    echo "> deploy: start"
    mkdir -p "$deploy_dir"
    wait_for_process_briefly "$binary_name"
    mv -f "./bin/${binary_name}" "$deploy_dir"
    echo "> deploy: success"
}

start() {
    echo "> run: start"
    mkdir -p "$logs_dir"
    binary_path="${deploy_dir}/${binary_name}"
    sudo setcap cap_net_bind_service=+ep "$binary_path"
    nohup "$binary_path address=':80'" &>> "${logs_dir}/${binary_name}-run.log" &
    echo "> run: success"
}

wait_for_process_briefly(){
    local pid=$(pidof "$@" || "")
    if [ -z "$pid" ]; then return; fi;
    local d=900
    echo "> deploy: waiting for previous shutdown.. (max:${$((d/10))}s)"
    local i=0
    while kill -0 "$pid"; do
        sleep 0.1s
        i=$((i+1))
        if [ $i -gt $d ]; then kill -9 "$pid" || true; fi;
    done
}

main $@