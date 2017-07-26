#!/usr/bin/env bash
## Author: Prasanna V. Loganathar

binary_name="nextfirst-core"
deploy_dir="${HOME}/workspace/bin/"
logs_dir="${HOME}/workspace/logs/"

main() {
    trap "echo '> script: incomplete termination requested'" TERM   
    set -e
    local dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )/../" && pwd )"    
    cd ${dir}
    build
    test_build
    deploy
    cd ${deploy_dir}
    start
}

build() {
    echo "> build: start"
    # go get -v -u github.com/golang/dep/cmd/dep || true
    # dep ensure || true
    go get -d -v || true
    go build -o "./bin/${binary_name}"
    echo "> build: done"
}

test_build() {
    echo "> test: start"
    go test
    echo "> test: done"
}

deploy() {
    echo "> deploy: start"
    mkdir -p "$deploy_dir"
    graceful_exit_or_kill "$binary_name" 90
    mv -f "./bin/${binary_name}" "$deploy_dir"
    echo "> deploy: done"
}

start() {
    echo "> run: start"
    mkdir -p "$logs_dir"
    binary_path="${deploy_dir}/${binary_name}"
    sudo setcap cap_net_bind_service=+ep "$binary_path"
    local log_file="${logs_dir}/${binary_name}-run.log"
    nohup "$binary_path" -address=":80" &>> $log_file &
    echo "> run: done"
}

graceful_exit_or_kill() {
    local pid=$(pidof "$1" || false)
    if [ -z "$pid" ]; then return; fi;
    local d=$(($2*10))
    echo "> deploy: waiting for previous shutdown.. (max: ${2}s)"
    local i=0
    while kill "$pid" >> /dev/null; do
        sleep 0.1s
        i=$((i+1))
        if [ $i -gt $d ]; then 
            echo "> deploy: forceful termation"
            kill -9 "$pid" >> /dev/null || true;
            break
        fi;
    done
}

main $@