#!/usr/bin/env bash
## Author: Prasanna V. Loganathar

set -e
dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )/../" && pwd )"
cd ${dir}

binary_name="nextfirst-core"
dest_dir="${HOME}/workspace/bin/"
logs_dir="${HOME}/workspace/logs/"

# Build
set +e
go get -u github.com/golang/dep/cmd/dep && dep ensure
go build -o ./bin/nextfirst-core

ret_val=$?
if [ $ret_val -ne 0 ]; then
    exit $ret_val
fi

# Deploy
set -e
mkdir -p "$dest_dir"
mkdir -p "$logs_dir"

wait_for_process_briefly(){
    local pid=$(pidof "$@")
    if [ -z "$pid" ]; then return; fi;
    local i=0
    while kill -0 "$pid"; do
        sleep 0.1s
        i=$((i+1))
        if [ $i -gt 900 ]; then kill -9 "$pid"; fi;
    done
}

wait_for_process_briefly "$binary_name"
mv -f "./bin/${binary_name}" "$dest_dir"

# Run
cd ${dest_dir}
binary_path="${dest_dir}/${binary_name}"
sudo setcap cap_net_bind_service=+ep "$binary_path"
nohup "$binary_path" address=":80" &>> "${logs_dir}/${binary_name}-run.log" &