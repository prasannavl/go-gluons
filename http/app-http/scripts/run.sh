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
    setup_bin_name
    PACKAGE_NAME="${BIN_NAME/%.exe}"
    LOGS_DIR="./logs"
    ASSETS_DIR="./www"
    CERT_DIR_CACHE="./cert-cache"
}

start() {
    graceful_exit_or_kill "./$BIN_NAME" 90
    echo "> run: start"
    mkdir -p "$LOGS_DIR" "$CERT_DIR_CACHE"
    local binary_path="./${BIN_NAME}"
    if [[ $(uname -s) != MINGW* ]]; then
        sudo setcap cap_net_bind_service=+ep "$binary_path"
    fi;
    local log_file_exec="${LOGS_DIR}/${PACKAGE_NAME}-exec.log"
    local log_file="${LOGS_DIR}/${PACKAGE_NAME}.log"  

    local cmd=$(echo "$binary_path" --address=":443" --root="${ASSETS_DIR}" --redirector=":80" --cert-dir="${CERT_DIR_CACHE}" --log="${log_file}")

    echo "cmd: " $cmd
    nohup $cmd &>> "${log_file_exec}" &

    echo "> run: done"
}

graceful_exit_or_kill() {
    local pid=$(pidof "$1" || false)
    if [ -z "$pid" ]; then return; fi;
    local d=$(($2*10))
    echo "> killer: waiting for previous shutdown.. (max: ${2}s)"
    local i=0
    while kill "$pid" &>> /dev/null; do
        sleep 0.1s
        i=$((i+1))
        if [ $i -gt $d ]; then 
            echo "> killer: forceful termination"
            kill -9 "$pid" &>> /dev/null || true;
            break
        fi;
    done
}

setup_bin_name() {
    if [ -n "$BIN_NAME" ]; then
        return
    fi;
    local cname=$(basename $(pwd))
    local cbin_name="${cname}"
    if [[ $(uname -s) == MINGW* ]]; then
        cbin_name="${cbin_name}.exe"
    fi
    BIN_NAME="${cbin_name}"
}

if [ -z "$@" ]; then
    run
else 
    "$@"
fi;
popd > /dev/null