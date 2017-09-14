#!/usr/bin/env bash
## Author: Prasanna V. Loganathar

binary_name="apicore"

app_dir="${HOME}/run/${binary_name}"
bin_dir="${app_dir}/bin"
assets_dir="${app_dir}/www"
logs_dir="${app_dir}/logs"
cert_cache_dir="${app_dir}/cert-cache"

build_assets_dir="./www"
build_target="./bin/${binary_name}"

main() {
    trap "echo '> script: incomplete termination requested'" TERM   
    set -e
    if [[ $# -gt 0 ]]; then
        eval "${@:1}"
        exit 0;
    fi;
    # Go to the parent directory of this script's dir
    local dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )/../" && pwd )"
    cd ${dir}
    build
    test_build
    deploy
    cd ${bin_dir}
    start
}

build() {
    echo "> build: start"
    echo "> build: updating deps"    
    go get -v -u github.com/golang/dep/cmd/dep || true
    dep ensure || true
    # go get -d -u ... || true
    echo "> build: compiling"
    go build -o "${build_target}"
    echo "> build: done"
}

test_build() {
    echo "> test: start"
    go test
    echo "> test: done"
}

deploy() {
    echo "> deploy: start"
    rm -rf "${assets_dir}"
    mkdir -p "$bin_dir" "$assets_dir"
    graceful_exit_or_kill "$binary_name" 90
    mv -f "${build_target}" "$bin_dir"
    cp -a $build_assets_dir/* "${assets_dir}/"
    echo "> deploy: done"
}

start() {
    graceful_exit_or_kill "$binary_name" 90    
    echo "> run: start"
    mkdir -p "$logs_dir" "$cert_cache_dir"
    binary_path="${bin_dir}/${binary_name}"
    sudo setcap cap_net_bind_service=+ep "$binary_path"
    local log_file_exec="${logs_dir}/${binary_name}-exec.log"
    local log_file="${logs_dir}/${binary_name}.log"
    nohup "$binary_path" --address=":443" --root="${assets_dir}" --redirector=":80" --cert-dir="${cert_cache_dir}" --log="${log_file}" &>> "${log_file_exec}" &
    echo "> run: done"
}

graceful_exit_or_kill() {
    local pid=$(pidof "$1" || false)
    if [ -z "$pid" ]; then return; fi;
    local d=$(($2*10))
    echo "> deploy: waiting for previous shutdown.. (max: ${2}s)"
    local i=0
    while kill "$pid" &>> /dev/null; do
        sleep 0.1s
        i=$((i+1))
        if [ $i -gt $d ]; then 
            echo "> deploy: forceful termination"
            kill -9 "$pid" &>> /dev/null || true;
            break
        fi;
    done
}

main "$@"