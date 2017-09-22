#!/usr/bin/env bash
## Author: Prasanna V. Loganathar

run() {
    init
    start
}

init() {
    if [ -z "$__INIT" ]; then 
        __INIT="1"
    else
        return 0
    fi
    trap "echo '> script: incomplete termination requested'" TERM   
    set -e
    # script dir
    local dir=$( dirname "${BASH_SOURCE[0]}" )
    # Go to the parent directory of this script's dir
    local pdir="$( cd "${dir}/../" && pwd )"
    pushd . > /dev/null
    cd ${pdir}
    setup_vars
}

setup_vars() {
    PKG_BASE_NAME="$(basename $(pwd))"
    BIN_NAME="${BIN_NAME:-${PKG_BASE_NAME}}"
    LOGS_DIR="./logs"
    # Pointing to a web root, that could be anywhere.
    # This dir is not touched by the deployment.
    WEBROOT_DIR="./www"
    CERT_CACHE_DIR="./cert-cache"
}

start() {
    stop
    echo "> run: start"
    mkdir -p "$LOGS_DIR" "$CERT_CACHE_DIR"
    local binary_path="./${BIN_NAME}"
    sudo setcap cap_net_bind_service=+ep "$binary_path"
    local log_file_exec="${LOGS_DIR}/${PKG_BASE_NAME}-exec.log"
    local log_file="${LOGS_DIR}/${PKG_BASE_NAME}.log"  

    local cmd=$(echo "$binary_path" --address=":443" --root="${WEBROOT_DIR}" --redirector=":80" --dapi-address="localhost:7000" --cert-dir="${CERT_CACHE_DIR}" --log="${log_file}")

    echo "cmd: " $cmd
    nohup $cmd &>> "${log_file_exec}" &

    echo "> run: done"
}

stop() {
    try_exit_or_kill "$BIN_NAME" 90
}

try_exit_or_kill() {
    local pid=$(pidof "$1" || false)
    if [ -z "$pid" ]; then return; fi;
    local d=$(($2*10))
    echo "> run: stopping: $1 (max-wait: ${2}s)"
    local i=0
    while kill "$pid" &>> /dev/null; do
        sleep 0.1s
        i=$((i+1))
        if [ $i -gt $d ]; then 
            echo "> run: sending SIG_KILL: $1"
            kill -9 "$pid" &>> /dev/null || true;
            break
        fi;
    done
    echo "> run: stopped: $1"
}

if [ -z "$@" ]; then
    run
else 
    "$@"
fi;
popd > /dev/null